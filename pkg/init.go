package pkg

import (
    "fmt"
    "io/ioutil"
	"os"
)

func Init() {
	fmt.Println("Initializing empty Goo site.")
	f_err := ioutil.WriteFile("site.yaml", []byte(initSite), 0644)
	Check(f_err)
	vd_err := os.Mkdir("views", 0755)
	Check(vd_err)
	pd_err := os.Mkdir("views/partials", 0755)
	Check(pd_err)
	dd_err := os.Mkdir("data", 0755)
	Check(dd_err)
	cd_err := os.Mkdir("data/posts", 0755)
	Check(cd_err)
	sd_err := os.Mkdir("static", 0755)
	Check(sd_err)
}

const initSite string = `
# See https://github.com/n-young/goo for documentation!

# Specify site name, output folder, static folder.
name: My Goo Site
output: build
static_src: static
static_dest: static

# Global variable files.
# global:
# 	everywhere: data/global.yaml


# Partials.
# partials:
#	header: views/partials/header.tmpl
#	footer: views/partials/footer.tmpl

# Pages.
# pages:
#     - title: Home
#       path: /
#       template: views/home.tmpl

# Collections.
# collections:
#     - title: Posts
#       path: /posts
#       template: views/post.tmpl
#       posts: data/posts
`