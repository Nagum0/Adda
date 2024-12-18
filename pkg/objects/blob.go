package objects

import (
	"adda/pkg/errors"
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
)

// Represents a blob object.
type Blob struct {
    FilePath string
    contents []byte
    hash     string
    Type     FileType
    Length   int
}

// Creates a new empty blob object.
func NewBlob(filePath string, fileType FileType) *Blob {
    return &Blob { 
        FilePath: filePath,
        contents: nil,
        hash: "",
        Type: fileType,
        Length: 0,
    }
}

// Returns the SHA-1 hash for the given file.
func (b *Blob) Hash() (string, error) {
    if b.hash != "" {
        return b.hash, nil
    }
    
    fileContents, err := os.ReadFile(b.FilePath)
    if err != nil {
        return "", errors.NewBlobError(fmt.Sprintf("Error while reading file for hash creation: %v", b.FilePath))
    }
    
    hasher := sha1.New()
    _, err = hasher.Write(fileContents)
    if err != nil {
        return "", errors.NewBlobError(fmt.Sprintf("Error while writing hash for file: %v", b.FilePath))
    }
    hashString := hex.EncodeToString(hasher.Sum(nil))
    b.hash = hashString

    return hashString, nil
}

// Returns the compressed contents of the blob object.
func (b *Blob) Contents() ([]byte, error) {
    if b.contents != nil {
        return b.contents, nil
    }

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
    b.contents = bytes

    return bytes, nil
}
