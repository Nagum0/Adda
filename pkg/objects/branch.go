package objects

import "os"

func CreateReferenceFile(refName string) error {
    path := ".adda/refs/heads/" + refName
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    return nil
}
