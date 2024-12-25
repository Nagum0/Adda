# ADDA

- A simple version control system built in go.
- Built for recreational purposes and learning.

## Commands

### Init

- Intializes an Adda repository.
- Usage:
``` bash
adda init
```
- Details:
  - Creates the .adda/ directory and the objects/, branches/, refs/, refs/heads/, subdirectories and the INDEX and HEAD files.

### Add

- Adds a file to the stagin area.
- Usage:
``` bash
adda add <filepath>
```
- Details
  - Creates a blob object for the file and stores it in the object database. The contents of the file are compressed using zlib.
  A hash is created based on the files contents using SHA-1 hashing. The blob is stored in the database in the directory with it's hashes prefix
  and the filename is the rest of the hash. After the creation of the blob object the INDEX file is updated to contain the the file path, hash and file type.

### Commit

- Commits the staged files and updates the HEAD of the branch.
- Usage:
``` bash
adda commit <message>
```
- Details:
  - x
