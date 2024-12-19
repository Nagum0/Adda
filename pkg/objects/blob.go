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

// Internally sets the SHA-1 hash for the blob object and returns the hash.
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

// Internally sets the contents of the blob object if we call this for the first time and 
// returns the compressed contents.
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

// Writes the blob object to the filesystem.
func (b *Blob) WriteBlob() error {
    hash, err := b.Hash()
    if err != nil {
        return err
    }

    // Creating the hash directory
    hashDir := hash[:2]
    err = os.Mkdir(".adda/objects/" + hashDir, os.ModePerm)
    if err != nil {
        return errors.NewBlobError(fmt.Sprintf("Error while creating hash directory: %v", err.Error()))
    }

    // Writing blob to a file
    hashFile := hash[2:]
    file, err := os.Create(".adda/objects/" + hashDir + "/" + hashFile)
    if err != nil {
        return errors.NewAddError(fmt.Sprintf("Error while creating blob file: %v", err.Error()), b.FilePath)
    }
    defer file.Close()

    contents, err := b.Contents()
    if err != nil {
        return err
    }

    file.Write(contents)

    return nil
}
