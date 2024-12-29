package objects

import (
	"adda/pkg/errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// An entry in the INDEX.
type Entry struct {
    FilePath string
    Hash     string
    FileType FileType
}

func (e Entry) String() string {
    return fmt.Sprintf("%v  %v\t%v", e.FileType, e.Hash, e.FilePath)
}

// The INDEX file.
// INDEX file format: filetype hash filepath
type Index struct {
    // Maps the file paths to the index entries.
    Entries map[string]Entry
}

func (i Index) String() string {
    indexString := ""
    
    for _, val := range i.Entries {
        indexString += val.String() + "\n"
    }

    return indexString
}

// Parses the INDEX file and retuns a pointer to the Index object.
func ParseIndex() (*Index, error) {
    indexFile, err := os.ReadFile(".adda/INDEX")
    if err != nil {
        return nil, errors.NewIndexError("Error while reading the INDEX file.")
    }
       
    indexFileString := string(indexFile)
    
    if len(indexFileString) == 0 {
        return &Index{ make(map[string]Entry) }, nil
    }

    indexFileLines := strings.Split(indexFileString, "\n")
    indexFileLines = indexFileLines[:len(indexFileLines ) - 1]
    indexObject := Index{ make(map[string]Entry) }
    
    for _, entry := range indexFileLines {
        splitEntry := strings.Fields(entry)

        if len(splitEntry) != 3 {
            return nil, errors.NewIndexError(fmt.Sprintf("Corrupted INDEX file: %v", splitEntry))
        }
        
        filePath := splitEntry[2]
        fileType, _ := strconv.Atoi(splitEntry[0])
        entryObject := Entry {
            FilePath: filePath,
            Hash: splitEntry[1],
            FileType: FileType(fileType),
        }

        indexObject.Entries[filePath] = entryObject
    }

    return &indexObject, nil
}

// Updates the INDEX file. If the file path already exists in the INDEX 
// it will only change the hash of the file.
func (index *Index) Update(blob Blob) error {
    // When calling this function the hash should already be generated
    hash, _ := blob.Hash()
    
    index.Entries[blob.FilePath] = Entry {
        FilePath: blob.FilePath,
        Hash: hash,
        FileType: blob.Type,
    }
    
    indexFile, err := os.Create(".adda/INDEX")
    if err != nil {
        return errors.NewIndexError("Error while updating the INDEX.")
    }
    defer indexFile.Close()
    indexFile.Write([]byte(index.String()))

    return nil
}

// Tries to get the hash of the given file path. If the file is not included in the index
// it will return an IndexError.
func (index Index) GetBlobHash(filePath string) (string, error) {
    entry, ok := index.Entries[filePath]
    if !ok {
        return "", errors.NewIndexError(fmt.Sprintf("File path: %v is not found in the INDEX", filePath))
    }

    return entry.Hash, nil
}
