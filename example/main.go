package main

import (
    "net/http"

    artefak "go-framework"
)

func main() {
    app := artefak.New()

    app.GET("/", func(c *artefak.Ctx) {
        c.String(http.StatusOK, "Halo")
    })

    app.Run(":8000")
}
