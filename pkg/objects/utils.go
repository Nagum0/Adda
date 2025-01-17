package objects

import "os"

// Checks whether the given hash already exists in the object database.
func HashExists(hash string) bool {
    hashPath := ".adda/objects/" + hash[:2] + "/" + hash[2:]

    _, err := os.Stat(hashPath)
    if os.IsNotExist(err) {
        return false
    }

    return true
}
