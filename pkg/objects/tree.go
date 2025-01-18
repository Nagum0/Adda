package objects

import (
	"adda/pkg/db"
	"adda/pkg/errors"
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"os"
	"sort"
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

// Takes a snapshot of the current directory structure of the staged files and retuns a Snaphot.
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

// Writes the snapshot's tree objects to the object database (zlib compression).
func (s Snapshot) WriteSnapshotToDatabase() error {
    for _, treeObject := range s {
        if db.HashExists(treeObject.Hash) {
            continue
        }
        
        treeString := treeObject.String()
        buffer := bytes.Buffer{}
        writer := zlib.NewWriter(&buffer)
        if _, err := writer.Write([]byte(treeString)); err != nil {
            return errors.NewCommitError("Error while writing the snapshot to the object database.")
        }
        writer.Close()

        hashPrefix := treeObject.Hash[:2]
        hashDirPath := ".adda/objects/" + hashPrefix + "/"
        _, err := os.Stat(hashDirPath)
        if os.IsNotExist(err) {
            if err := os.Mkdir(hashDirPath, os.ModePerm); err != nil {
                return errors.NewCommitError(fmt.Sprintf("Error while creating hash directory for tree object: %v", treeObject.DirName))
            }
        }

        hashFilePath := ".adda/objects/" + hashPrefix + "/" + treeObject.Hash[2:]
        file, err := os.Create(hashFilePath)
        if err != nil {
            return errors.NewCommitError(fmt.Sprintf("Error while creating tree object file for tree object: %v", treeObject.DirName))
        }
        defer file.Close()
        file.Write(buffer.Bytes())
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

    tree.Hash = fmt.Sprintf("%x", sha1.Sum([]byte(tree.String())))
    snapshot[dirName] = tree
}
