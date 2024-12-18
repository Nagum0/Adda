package main

import (
	"adda/pkg/commands"
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
        err := commands.Init()
        if err != nil {
            fmt.Println(err.Error())
            os.Exit(1)
        }

        break
    case "add":
        err := commands.Add(args[1])
        if err != nil {
            fmt.Println(err.Error())
            os.Exit(1)
        }

        break
    case "commit":
        fmt.Println("COMMIT")
        break
    default:
        fmt.Println("Unknown command: ", args[0])
        break
    }
}
