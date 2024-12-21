package index

import (
	"adda/pkg/objects"
	"fmt"
)

// A file in the INDEX.
type Entry struct {
    FileName string
    Hash     string
    FileType objects.FileType
}

func (e Entry) String() string {
    return fmt.Sprintf("%v  %v\t%v", e.FileType, e.Hash, e.FileName)
}

// The INDEX file.
type Index struct {
    // Maps the file paths to the index entries.
    Entries map[string]Entry
}

// Parses the INDEX file and retuns a pointer to the Index object.
func ParseIndex() (*Index, error) {
    // indexFile, err := os.ReadFile(".adda/INDEX")
    // if err != nil {
    //     return nil, nil
    // }

    return nil, nil
}

// Updates the INDEX file. If the file path already exists in the INDEX 
// it will only change the hash of the file.
func (index *Index) Update(blob objects.Blob) error {
    return nil
}
