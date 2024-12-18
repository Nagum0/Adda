package objects

import (
	"adda/pkg/errors"
	"bytes"
	"compress/zlib"
	"fmt"
	"os"
)

// Represents a blob object.
type Blob struct {
    FilePath string
    Contents []byte
    Hash     string
    Type     FileType
    Length   int
}

func NewBlob(filePath string, fileType FileType) *Blob {
    return &Blob { 
        FilePath: filePath,
        Contents: make([]byte, 0),
        Hash: "",
        Type: fileType,
        Length: 0,
    }
}

func (b *Blob) GenerateHash() string {
    return ""
}

func (b *Blob) GenerateContents() ([]byte, error) {
    fileContents, err := os.ReadFile(b.FilePath)
    if err != nil {
        return nil, errors.NewBlobError(fmt.Sprintf("Error while reading file for blob object creation: %v", b.FilePath))
    }

    b.Length = len(fileContents)
    buffer := bytes.Buffer{}
    writer := zlib.NewWriter(&buffer)
    _, err = writer.Write(fileContents)
    if err != nil {
        return nil, errors.NewBlobError(fmt.Sprintf("Error while creating blob object for file: %v", b.FilePath))
    }
    writer.Close()

    bytes := buffer.Bytes()
    b.Contents = bytes

    return bytes, nil
}
