package main

import (
	"fmt"
	"os"

	goo "github.com/n-young/goo/pkg"
)

// Main entry.
func main() {
	// Parse args; must have at least a command.
	args := os.Args;
	if len(args) < 2 {
		fmt.Println("Usage: goo <command>");
		return;
	}
	
	// Switch on the command.
	switch args[1] {
		// Help - prints out the help menu.
		case "help": goo.Help();
		// Init - initializes and empty goo site in the current directory.
		case "init": goo.Init();
		// Build - given a <site.yml>, builds the goo site.
		case "build": {
			if len(args) != 3 {
				fmt.Println("Usage: goo build <site.yml>");
				return;
			}
			goo.Build(args[2]);
		}
		// Watch - repeatedly runs Build whenever a file changes.
		case "watch": {
			if len(args) != 3 {
				fmt.Println("Usage: goo watch <site.yml>");
				return;
			}
			goo.Watch(args[2]);
		}
		// Default: print out usage details.
		default: {
			fmt.Println("Usage: goo <command>");
		}
	}
}
