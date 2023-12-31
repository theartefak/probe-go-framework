package artefak

import (
    "fmt"
    "log"
    "net/http"
    "runtime"
    "strings"
)

func trace(message string) string {
    var pcs [32]uintptr
    var str strings.Builder

    n := runtime.Callers(3, pcs[:])
    str.WriteString(message + "\nTraceback:")

    for _, pc := range pcs[:n] {
        fn := runtime.FuncForPC(pc)
        file, line := fn.FileLine(pc)
        str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
    }

    return str.String()
}

func Recovery() HandlerFunc {
    return func(c *Ctx) {
        defer func() {
            if err := recover(); err != nil {
                message := fmt.Sprintf("%s", err)
                log.Printf("%s\n\n", trace(message))
                c.Fail(http.StatusInternalServerError, "Internal Server Error")
            }
        }()

        c.Next()
    }
}
