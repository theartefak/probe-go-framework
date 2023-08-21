package artefak

import (
    "log"
    "time"
)

func Logger() HandlerFunc {
    return func(c *Ctx) {
        t := time.Now()
        c.Next()
        log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
    }
}
