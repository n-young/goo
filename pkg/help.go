package pkg

import "fmt"

func Help() {
	fmt.Println(
		`Usage: goo <command>
Compile a Goo site.

Available commands:
  help: Bring up this text.
  init: Create a blank Goo site.
  build <site.yml>: Compile a Goo site.
  serve <site.yml> <port> [--nobuild]: Hot-reloading compiles.`)
}
