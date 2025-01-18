package db

import (
	"adda/pkg"
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

// ---------------------------------------------------------------------------------------
//                                  OBJECT DATABASE FUNCTIONALITY
// ---------------------------------------------------------------------------------------

// Types that implement this interface can write to the object database using 
// the DBWrite() function. It is expected that the types hold the hash and contents 
// that are needed to be written to the database.
type DBWriter interface {
    DBWrite() error
}

// Reads the database at the given hash.
func DBRead(hash string) ([]byte, error) {
    hashPath := pkg.OBJECTS_PATH + hash[:2] + "/" + hash[2:]
    contents, err := os.ReadFile(hashPath)
    if err != nil {
        return nil, err
    }

    return contents, nil
}

// Checks whether the given hash already exists in the object database.
func HashExists(hash string) bool {
    hashPath := pkg.OBJECTS_PATH + hash[:2] + "/" + hash[2:]

    _, err := os.Stat(hashPath)
    if os.IsNotExist(err) {
        return false
    }

    return true
}

// Generate a SHA-1 hash for the given byte slice. The hash will be returned as a string.
func GenSHA1(contents []byte) string {
    return fmt.Sprintf("%x", sha1.Sum(contents))
}

// Compresses the given byte slice using zlib.
func ZlibCompress(contents []byte) []byte {
    buffer := bytes.Buffer{}
    zlibWriter := zlib.NewWriter(&buffer)
    zlibWriter.Write(contents)
    zlibWriter.Close()
    return buffer.Bytes()
}

// Compresses the given string using zlib.
func ZlibCompressString(s string) []byte {
    return ZlibCompress([]byte(s))
}

// Decompress the byte slice from zlib.
func Decompress(contents []byte) ([]byte, error) {
    bytes := bytes.NewBuffer(contents)
    reader, err := zlib.NewReader(bytes)
    if err != nil {
        return nil, err
    }
    defer reader.Close()
    
    bytes.Reset()
    if _, err := io.Copy(bytes, reader); err != nil {
        return nil, err
    }

    return bytes.Bytes(), nil
}

// Decompress the byte slice from zlib compression to a string.
func DecompressToString(contents []byte) (string, error) {
    b, err := Decompress(contents)
    return string(b), err
}

// ---------------------------------------------------------------------------------------
//                                  REFS FUNCTIONALITY
// ---------------------------------------------------------------------------------------

// Reads the contents of the HEAD file.
func ReadHEAD() (string, error) {
    head, err := os.ReadFile(pkg.HEAD_PATH)
    if err != nil {
        return "", err
    }

    return string(head), nil
}

// Sets the HEAD file to the given reference.
func SetHEAD(ref string) error {
    head, err := os.Create(pkg.HEAD_PATH)
    if err != nil {
        return err
    }
    defer head.Close()

    if _, err := head.WriteString(ref); err != nil {
        return err
    }

    return nil
}

// Reads the contents of the given branch's reference head.
func ReadRefHead(branchName string) (string, error) {
    b, err := os.ReadFile(pkg.REFS_HEADS_PATH + branchName)
    if err != nil {
        return "", err
    }

    return string(b), nil
}

// Sets the given branch's reference head to the given hash.
// If the head file for the branch doesn't exist it creates it at .adda/refs/heads/branchName.
func SetRefsHead(branchName string, hash string) error {
    headFile, err := os.Create(pkg.REFS_HEADS_PATH + branchName)
    if err != nil {
        return err
    }
    defer headFile.Close()

    if _, err = headFile.WriteString(hash); err != nil {
        return err
    }

    return nil
}
