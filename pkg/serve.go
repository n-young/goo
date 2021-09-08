package pkg

import (
	"fmt"
	"net/http"
)

func Serve(file string, port string) {
	config := ParseConfig(file)
	directory := config.Output
	http.Handle("/", http.FileServer(http.Dir(directory)))
	fmt.Printf("Serving %s on HTTP port: %s\n", directory, port)
	http.ListenAndServe(":" + port, nil)
}