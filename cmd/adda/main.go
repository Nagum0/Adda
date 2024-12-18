package main

import (
	"adda/pkg/adda"
	"fmt"
	"os"
)

func main() {
    args := os.Args[1:]

    if len(args) == 0 {
        fmt.Println("Usage: ...")
        os.Exit(1)
    }

    switch args[0] {
    case "init":
        err := adda.Init()
        if err != nil {
            fmt.Println(err.Error())
            os.Exit(1)
        }
        break
    case "add":
        fmt.Println("ADD")
        break
    case "commit":
        fmt.Println("COMMIT")
        break
    default:
        fmt.Println("Unknown command: ", args[0])
        break
    }
}
