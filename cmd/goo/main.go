package main

import (
	"flag"
	"fmt"
	// "os"
	"strconv"

	goo "github.com/n-young/goo/pkg"
)

// Main entry.
func main() {
	// cli flags
	nobuild := flag.Bool("nobuild", false, "set this flag if you don't want to build using serve")
	flag.Parse()

	// Parse args; must have at least a command.
	args := flag.Args();
	if len(args) < 1 {
		fmt.Println("Usage: goo <command>");
		return;
	}
	
	// Switch on the command.
	switch args[0] {
		// Help - prints out the help menu.
		case "help": goo.Help();
		// Init - initializes and empty goo site in the current directory.
		case "init": goo.Init();
		// Build - given a <site.yml>, builds the goo site.
		case "build": {
			if len(args) != 2 {
				fmt.Println("Usage: goo build <site.yml>");
				return;
			}
			goo.Build(args[1]);
		}
		// Serve - .
		case "serve": {
			if len(args) != 3 && len(args) != 4 {
				fmt.Println("Usage: goo serve <site.yml> <port>");
				return;
			} else if _, err := strconv.Atoi(args[2]); err != nil {
				fmt.Println("invalid port")
				return
			}
			if !*nobuild {
				goo.Build(args[1]);
			}
			goo.Serve(args[1], args[2]);
		}
		// Default: print out usage details.
		default: {
			fmt.Println("Usage: goo <command>");
		}
	}
}
