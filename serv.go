package main

import (
	"fmt"
	"net/http"

	"github.com/theartefak/trit/artefak"
)

func main() {
    app := artefak.New()

    app.GET("/", func(w http.ResponseWriter, r *http.Request) {
    	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
    })

    app.Run(":8000")
}
