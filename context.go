package artefak

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type H map[string]interface{}

type Ctx struct {
    Writer     http.ResponseWriter
    Req        *http.Request
    Path       string
    Method     string
    Params     map[string]string
    StatusCode int
    handlers   []HandlerFunc
    index      int
    artefak    *Artefak
}

func NewCtx(w http.ResponseWriter, req *http.Request) *Ctx {
    return &Ctx{
        Path   : req.URL.Path,
        Method : req.Method,
        Req    : req,
        Writer : w,
        index  : -1,
    }
}

func (c *Ctx) Next() {
    c.index++
    s := len(c.handlers)
    for ; c.index < s; c.index++ {
        c.handlers[c.index](c)
    }
}

func (c *Ctx) Fail(code int, err string) {
    c.index = len(c.handlers)
    c.JSON(code, H{"message": err})
}

func (c *Ctx) Param(key string) string {
    value, _ := c.Params[key]
    return value
}

func (c *Ctx) PostForm(key string) string {
    return c.Req.FormValue(key)
}

func (c *Ctx) Query(key string) string {
    return c.Req.URL.Query().Get(key)
}

func (c *Ctx) Status(code int) {
    c.StatusCode = code
    c.Writer.WriteHeader(code)
}

func (c *Ctx) SetHeader(key string, value string) {
    c.Writer.Header().Set(key, value)
}

func (c *Ctx) String(code int, format string, values ...interface{}) {
    c.SetHeader("Content-Type", "text/plain")
    c.Status(code)
    c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Ctx) JSON(code int, obj interface{}) {
    c.SetHeader("Content-Type", "application/json")
    c.Status(code)
    encoder := json.NewEncoder(c.Writer)
    if err := encoder.Encode(obj); err != nil {
        http.Error(c.Writer, err.Error(), 500)
    }
}

func (c *Ctx) Data(code int, data []byte) {
    c.Status(code)
    c.Writer.Write(data)
}

func (c *Ctx) HTML(code int, name string, data interface{}) {
    c.SetHeader("Content-Type", "text/html")
    c.Status(code)

    if err := c.artefak.htmlTemplates.ExecuteTemplate(c.Writer, name+".html", data); err != nil {
        c.Fail(500, err.Error())
    }
}
