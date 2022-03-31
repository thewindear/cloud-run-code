package ctl

import (
	"cloud-run-code/context"
	"cloud-run-code/service"
	context2 "context"
	"io/ioutil"
	"net/http"
	"time"
)

// RunController Post /run?lang={lang}
func RunController(writer http.ResponseWriter, request *http.Request) {
	ctx := &context.Context{Writer: writer, Req: request}
	if ctx.Req.Method != "POST" {
		ctx.NotAllow()
		return
	}
	lang := ctx.Get("lang", "")
	if lang == "" {
		ctx.Ret(http.StatusBadRequest, lang, "lang parameter can't be empty")
		return
	} else {
		if !service.DockerRunner.RunnerExists(lang) {
			ctx.Ret(http.StatusBadRequest, lang, lang+"not support")
			return
		}
		body, _ := ioutil.ReadAll(request.Body)
		if len(body) == 0 {
			ctx.Ret(http.StatusBadRequest, lang, "The code required to be executed cannot be empty")
			return
		}

		code := string(body)
		cancelCtx, cancelFn := context2.WithTimeout(context2.Background(), time.Second*time.Duration(service.DockerRunner.Timeout))
		defer cancelFn()

		content, err := service.DockerRunner.Exec(cancelCtx, lang, code)
		if err == nil {
			ctx.OK(lang, string(content), "execute success")
			return
		}
		if err == service.TimeoutError {
			ctx.Ret(http.StatusRequestTimeout, lang, "execute timeout")
			return
		}
		ctx.Ret(http.StatusInternalServerError, lang, string(content))
	}
}
