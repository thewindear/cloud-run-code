### Cloud-Run-Code

> 在线运行代码
> 服务依赖docker环境
> 
```shell
    go run main.go --port=端口号 --config=配置文件路径
```

#### 手动其他添加语言
编辑`docker_config.json`文件 添加对应 key value 即可

```json
{
  "docker_runner": {
    "go": {
      "ext": "go",
      "filename": "main.go",
      "image": "golang:1.18.0-alpine3.15",
      "cmd": "go run main.go"
    }  
  }  
}
```
`docker_runner`字段说明

| 字段名 | 说明 |
|-------|------|
| ext | 文件后续 |
| filename | 容器中运行的文件名 |
| image | 镜像名 |
| cmd | 启动容器后运行的命令 |

#### 接口说明

获取所有支持的语言列表

`GET /languages`
```json
[
  "go",
  "python",
  "php",
  "java",
  "node"
]
```

`POST /run?lang={lang}`
* 参数说明 lang : 为 languages 接口返回的列表
* 请求体：具体代码
```json
{
  "lang": "{lang}",
  "message": "execute success",
  "result": "hello world"
}
```