package main

import (
	"fmt"
	"os"

	goo "github.com/n-young/goo/pkg"
)

func main() {
	args := os.Args;
	if len(args) < 2 {
		fmt.Println("Usage: goo <command>");
		return;
	}
	
	switch args[1] {
		case "help": goo.Help();
		case "init": goo.Init();
		case "build": {
			if len(args) != 3 {
				fmt.Println("Usage: goo build <site.yml>");
				return;
			}
			goo.Build(args[2]);
		}
		case "watch": {
			if len(args) != 3 {
				fmt.Println("Usage: goo watch <site.yml>");
				return;
			}
			goo.Watch(args[2]);
		}
		default: {
			fmt.Println("Usage: goo <command>");
		}
	}
}
