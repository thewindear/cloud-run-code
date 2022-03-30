package ctl

import (
	"cloud-run-code/context"
	"cloud-run-code/service"
	context2 "context"
	"io/ioutil"
	"net/http"
	"time"
)

func RunController(writer http.ResponseWriter, request *http.Request) {
	ctx := &context.Context{Writer: writer, Req: request}
	if ctx.Req.Method != "POST" {
		ctx.NotAllow()
		return
	}
	lang := ctx.Get("lang", "")
	if lang == "" {
		ctx.Ret(http.StatusBadRequest, lang, "语言不能为空")
		return
	} else {
		if !service.DockerRunner.RunnerExists(lang) {
			ctx.Ret(http.StatusBadRequest, lang, "暂不支持")
			return
		}
		body, _ := ioutil.ReadAll(request.Body)
		if len(body) == 0 {
			ctx.Ret(http.StatusBadRequest, lang, "代码不能为空")
			return
		}

		code := string(body)
		cancelCtx, cancelFn := context2.WithTimeout(context2.Background(), time.Second*time.Duration(service.DockerRunner.Timeout))
		defer cancelFn()

		content, err := service.DockerRunner.Exec(cancelCtx, lang, code)
		if err == nil {
			ctx.OK(lang, string(content), "执行成功")
			return
		}
		if err == service.TimeoutError {
			ctx.Ret(http.StatusRequestTimeout, lang, "代码执行超时")
			return
		}
		ctx.Ret(http.StatusInternalServerError, lang, string(content))
	}
}
