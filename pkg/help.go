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
  watch <site.yml>: Hot-reloading compiles.`)
}
