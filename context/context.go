package context

import (
    "encoding/json"
    "net/http"
    "strings"
)

type Result struct {
    Lang    string `json:"lang,omitempty"`
    Message string `json:"message,omitempty"`
    Result  string `json:"result,omitempty"`
}

type Context struct {
    Writer http.ResponseWriter
    Req    *http.Request
}

func (c *Context) IsGet() bool {
    return c.Method() == http.MethodGet
}

func (c *Context) IsPost() bool {
    return c.Method() == http.MethodPost
}

func (c *Context) Method() string {
    return strings.ToUpper(c.Req.Method)
}

func (c *Context) NotAllow() {
    c.Writer.WriteHeader(http.StatusMethodNotAllowed)
}

func (c *Context) JSON(code int, data interface{}) error {
    c.Writer.Header().Add("Content-Type", "application/json")
    content, _ := json.Marshal(data)
    c.res(code, content)
    return nil
}

func (c *Context) RunOK(lang, result, message string) {
    _ = c.JSON(http.StatusOK, &Result{lang, message, result})
}

func (c *Context) Bad(message string) {
    c.RunRet(http.StatusBadRequest, message)
}

func (c *Context) Error(message string) {
    c.RunRet(http.StatusInternalServerError, message)
}

func (c *Context) Timeout(message string) {
    c.RunRet(http.StatusRequestTimeout, message)
}

func (c *Context) RunRet(code int, message string) {
    _ = c.JSON(code, &Result{"", message, ""})
}

func (c *Context) res(code int, data []byte) {
    c.Writer.WriteHeader(code)
    _, _ = c.Writer.Write(data)
}

func (c *Context) Get(key string, def string) string {
    val := c.Req.URL.Query().Get(key)
    if val == "" {
        return def
    }
    return val
}
