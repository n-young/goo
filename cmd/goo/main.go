package main

import (
	"os"
	"fmt"
	"strconv"

	goo "github.com/n-young/goo/pkg"
)

// Main entry.
func main() {
	// Parse args; must have at least a command.
	args := os.Args
	fmt.Println(args)
	if len(args) == 1 {
		fmt.Println("Usage: goo <command>")
		return
	}

	// Switch on the command.
	switch args[1] {
	// Help - prints out the help menu.
	case "help":
		goo.Help()
	// Init - initializes and empty goo site in the current directory.
	case "init":
		goo.Init()
	// Build - given a <site.yml>, builds the goo site.
	case "build":
		if len(args) != 3 {
			fmt.Println("Usage: goo build <site.yml>")
			return
		}
		goo.Build(args[2])
	// Serve - .
	case "serve":
		if len(args) != 4 && len(args) != 5 {
			fmt.Println("Usage: goo serve <site.yml> <port>")
			return
		} else if _, err := strconv.Atoi(args[2]); err != nil {
			fmt.Println("invalid port")
			return
		}
		goo.Build(args[2])
		goo.Serve(args[2], args[3])
	// Default: print out usage details.
	default:
		fmt.Println("Usage: goo <command>")
	}
}
