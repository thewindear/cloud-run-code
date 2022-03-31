package context

import (
    "encoding/json"
    "net/http"
)

type LangResult struct {
    Lang    string `json:"lang"`
    Message string `json:"message"`
    Result  string `json:"result"`
}

type Context struct {
    Writer http.ResponseWriter
    Req    *http.Request
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

func (c *Context) OK(lang, result, message string) {
    _ = c.JSON(http.StatusOK, &LangResult{lang, message, result})
}

func (c *Context) Ret(code int, lang, message string) {
    _ = c.JSON(code, &LangResult{lang, message, ""})
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
