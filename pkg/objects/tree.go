package objects

import (
	"adda/pkg"
	"adda/pkg/db"
	"adda/pkg/errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// A blob object in the tree object. Holds the hash and the original file name of the blob.
type TreeBlob struct {
    FileName string
    Hash     string
}

func NewTreeBlob(fileName string, hash string) *TreeBlob {
    return &TreeBlob {
        FileName: fileName,
        Hash: hash,
    }
}

// A tree object holds the blob objects of a directory and it's hash.
// Tree object formatting: <object type> <hash> <file name>. The object type can be 0 for a blob
// representing a file or 1 for a subtree representing a subdirectory.
type TreeObject struct {
    // Name of the directory the tree is representing.
    DirName string
    // The hash of the tree object.
    Hash    string
    // The tree's subdirectories
    SubDirs []*TreeObject
    // The tree's blob objects   
    Blobs   []TreeBlob
}

func NewTreeObject(dirName string) *TreeObject {
    return &TreeObject {
        DirName: dirName,
        SubDirs: make([]*TreeObject, 0),
        Blobs: make([]TreeBlob, 0),
    }
}

// Checks whether the tree object contains the given directory name as one of it's subdirectories.
func (tree TreeObject) ContainsSubDir(dirName string) bool {
    for _, subDir := range tree.SubDirs {
        if dirName == subDir.DirName {
            return true
        }
    }

    return false
}

// Checks whether the tree object contains the given file name as one of it's blob objects.
func (tree TreeObject) ContainsBlobFile(fileName string) bool {
    for _, blob := range tree.Blobs {
        if blob.FileName == fileName {
            return true
        }
    }

    return false
}

// Takes a string representation of a tree object and parses it into a TreeObject.
func ParseTreeObjectString(s string, isRoot bool) (*TreeObject, error) {
    treeObject := *NewTreeObject("")
    if isRoot {
        treeObject.DirName = "."
    }
    treeHash := db.GenSHA1([]byte(s))
    treeObject.Hash = treeHash

    lines := strings.Split(s, "\n")
    lines = lines[:len(lines) - 1]
    
    for _, line := range lines {
        lineFields := strings.Fields(line)

        if len(lineFields) != 3 {
            return nil, errors.NewTreeError("Illegal tree object string format.")
        }

        fileType, err := strconv.Atoi(lineFields[0])
        if err != nil {
            return nil, errors.NewTreeError(err.Error())
        }

        switch fileType {
        // Blob file
        case 0:
            blob := TreeBlob{
                Hash: lineFields[1],
                FileName: lineFields[2],
            }
            treeObject.Blobs = append(treeObject.Blobs, blob)
            break
        // Subtree
        case 1:
            subTreeHash := lineFields[1]
            subTreeName := lineFields[2]
            subTreeBytes, err := db.DBRead(subTreeHash)
            if err != nil {
                return nil, errors.NewTreeError(err.Error())
            }

            subTreeString, err := db.DecompressToString(subTreeBytes)
            if err != nil {
                return nil, errors.NewTreeError(err.Error())
            }

            subTree, err := ParseTreeObjectString(subTreeString, false)
            subTree.DirName = subTreeName
            treeObject.SubDirs = append(treeObject.SubDirs, subTree)
            break
        default:
            return nil, errors.NewTreeError(fmt.Sprintf("Unknown file type: %v", fileType))
        }
    }

    return &treeObject, nil
}

// Parse tree object into a string. Tree object format: <object type> <hash> <file/directory name>.
func (tree TreeObject) String() string {
    s := ""
    
    for _, blob := range tree.Blobs {
        s += fmt.Sprintf("0 %v\t%v\n", blob.Hash, blob.FileName)
    }

    for _, subDir := range tree.SubDirs {
        s += fmt.Sprintf("1 %v\t%v\n", subDir.Hash, subDir.DirName)
    }

    return s
}

// Snapshot of the directory structure at the time of a given commit. 
// Maps the directory names to tree objects.
type Snapshot map[string]*TreeObject

func NewSnapshot() *map[string]*TreeObject {
    return &map[string]*TreeObject{}
}

// Takes a snapshot of the current directory structure of the staged files and retuns a Snapshot.
func TakeSnapshot(indexFile Index) Snapshot {
    snapshot := *NewSnapshot()
    snapshot["."] = NewTreeObject(".")
    
    for _, entry := range indexFile.Entries {
        dirs := strings.Split(entry.FilePath, "/")
        currentDir := "."

        for i := 0; i < len(dirs) - 1; i++ {
            nextDir := dirs[i]
            nextDirTree := NewTreeObject(nextDir)

            if _, ok := snapshot[nextDir]; !ok {
                snapshot[nextDir] = nextDirTree
            }

            // Add the nextDir to the currentDir as a child if it doesn't already contain it
            if currentDirTree := snapshot[currentDir]; !currentDirTree.ContainsSubDir(nextDir) {
                currentDirTree.SubDirs = append(currentDirTree.SubDirs, nextDirTree)
                snapshot[currentDir] = currentDirTree
            }

            currentDir = nextDir
        }

        fileName := dirs[len(dirs) - 1]
        currentDirTree := snapshot[currentDir]
        currentDirTree.Blobs = append(currentDirTree.Blobs, *NewTreeBlob(fileName, entry.Hash))
        snapshot[currentDir] = currentDirTree
    }

    generateTreeHashes(snapshot, ".")
    
    return snapshot
}

// TODO: BuildIndex
// Takes a snapshot and builds an index file from it.
func (s Snapshot) BuildIndex() *Index {
    panic("todo")
}

// TODO: BuildBranchsLatestSnapshot
// Takes the hash of the root tree of the branchs latest commit object and builds the snapshot from it.
func BuildBranchsLatestSnapshot(rootTreeHash string) (*Snapshot, error) {
    snapshot := Snapshot{}
    rootTree, err := getRootTree(rootTreeHash)
    if err != nil {
        return nil, errors.NewBranchError(err.Error())
    }
    
    snapshot["."] = rootTree

    fmt.Println(rootTree)
    
    return &snapshot, nil
}

func getRootTree(rootTreeHash string) (*TreeObject, error) {
    rootTree, err := db.DBRead(rootTreeHash)
    if err != nil {
        return nil, errors.NewTreeError(err.Error())
    }
    
    rootTreeString, err := db.DecompressToString(rootTree)
    if err != nil {
        return nil, errors.NewTreeError(err.Error())
    }

    rootTreeObject, err := ParseTreeObjectString(rootTreeString, true)
    if err != nil {
        return nil, errors.NewTreeError(err.Error())
    }

    return rootTreeObject, nil
}

// Writes the snapshot's tree objects to the object database (zlib compression).
func (s Snapshot) DBWrite() error {
    for _, treeObject := range s {
        if db.HashExists(treeObject.Hash) {
            continue
        }
        
        compressedBytes := db.ZlibCompressString(treeObject.String())
        hashPrefix := treeObject.Hash[:2]
        hashDirPath := pkg.OBJECTS_PATH + hashPrefix + "/"
        _, err := os.Stat(hashDirPath)
        if os.IsNotExist(err) {
            if err := os.Mkdir(hashDirPath, os.ModePerm); err != nil {
                return errors.NewCommitError(fmt.Sprintf("Error while creating hash directory for tree object: %v", treeObject.DirName))
            }
        }

        hashFilePath := pkg.OBJECTS_PATH + hashPrefix + "/" + treeObject.Hash[2:]
        file, err := os.Create(hashFilePath)
        if err != nil {
            return errors.NewCommitError(fmt.Sprintf("Error while creating tree object file for tree object: %v", treeObject.DirName))
        }
        defer file.Close()
        file.Write(compressedBytes)
    }

    return nil
}

func (snap Snapshot) String() string {
    s := ""
    
    for dirName, treeObject := range snap {
        s += fmt.Sprintf("%v -> %v:\n%v\n", dirName, treeObject.Hash, treeObject)
    }

    return s
}

// Generate the hashes for the tree objects in the given snapshot.
func generateTreeHashes(snapshot Snapshot, dirName string) {
    for _, subDir := range snapshot[dirName].SubDirs {
        if subDir.Hash == "" {
            generateTreeHashes(snapshot, subDir.DirName)
        }
    }

    tree := snapshot[dirName]

    // Sort the tree blobs (by file name) and sub directories (by directory name) for consistent hashing
    sort.Slice(tree.Blobs, func(i, j int) bool {
        return tree.Blobs[i].FileName < tree.Blobs[j].FileName
    })

    sort.Slice(tree.SubDirs, func(i, j int) bool {
        return tree.SubDirs[i].DirName < tree.SubDirs[j].DirName
    })

    tree.Hash = db.GenSHA1([]byte(tree.String()))
    snapshot[dirName] = tree
}
