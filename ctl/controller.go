package ctl

import (
    "cloud-run-code/context"
    "cloud-run-code/service"
    context2 "context"
    "io/ioutil"
    "net/http"
    "time"
)

// LanguagesController Get /languages
func LanguagesController(writer http.ResponseWriter, request *http.Request) {
    ctx := &context.Context{Writer: writer, Req: request}
    if !ctx.IsGet() {
        ctx.NotAllow()
        return
    }
    var languages []string
    for lang := range service.DockerRunner.Runners {
        languages = append(languages, lang)
    }
    _ = ctx.JSON(http.StatusOK, languages)
    return
}

// RunController Post /run?lang={lang}
func RunController(writer http.ResponseWriter, request *http.Request) {
    ctx := &context.Context{Writer: writer, Req: request}
    if !ctx.IsPost() {
        ctx.NotAllow()
        return
    }
    lang := ctx.Get("lang", "")
    if lang == "" {
        ctx.Bad("lang parameter can't be empty")
        return
    } else {
        if !service.DockerRunner.RunnerExists(lang) {
            ctx.Bad(lang + "not support")
            return
        }
        body, _ := ioutil.ReadAll(request.Body)
        if len(body) == 0 {
            ctx.Bad("The code required to be executed cannot be empty")
            return
        }

        code := string(body)
        cancelCtx, cancelFn := context2.WithTimeout(context2.Background(), time.Second*time.Duration(service.DockerRunner.Timeout))
        defer cancelFn()

        content, err := service.DockerRunner.Exec(cancelCtx, lang, code)
        if err == nil {
            ctx.RunOK(lang, string(content), "execute success")
            return
        }
        if err == service.TimeoutError {
            ctx.Timeout("execute timeout")
            return
        }
        ctx.Error(string(content))
    }
}
