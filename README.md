# Pi

![build.yaml](https://github.com/go-laeo/pi/actions/workflows/build.yaml/badge.svg) [![Go Reference](https://pkg.go.dev/badge/github.com/go-laeo/pi.svg)](https://pkg.go.dev/github.com/go-laeo/pi)

[English](https://github.com/go-laeo/pi/blob/main/README_EN.md)

“Pi”是一个简洁轻量高性能的路由组件，通过使用泛型来减少后端接口开发中的模版代码。

# 代码示例

```Go
type UserData struct {
    Name string
    Password string
}

func h(ctx pi.Context) error {
    data := &UserData{}
    err := pi.Format(ctx, data)
    if err != nil {
        return pi.NewError(400, err.Error())
    }

    // do sth. actions...

    return ctx.Text("hello, world!")
}

sm := pi.NewServerMux(context.Background())
sm.Route("/api/v1/users").Post(h)

http.ListenAndServe("localhost:8080", sm)
```

# 安装

```shell
go get -u github.com/go-laeo/pi
```

# 特点
* [x] 基于前缀树的高性能路由功能，支持路由参数提取、通配路由等功能
* [x] ~~兼容 `net/http` (`pi.HandlerFunc` 实现了 `http.Handler`)~~
* [x] ~~Auto~~使用泛型函数 `pi.Format[T any]()` 来主动解析请求体
* [x] 路由中间件由 `pi.(ServerMux).Use()` 或 `pi.(HandlerFunc).Connect()` 进行注入
* [x] 内置针对 SPA 应用优化的 `pi.FileServer`
* [x] 无外部库依赖，无供应链攻击风险
* [x] 完备的单元测试与性能测试
# 更多示例

查看 `_examples` 目录。

# Related Projects
* [Wetalk](https://github.com/go-laeo/wetalk)
# 开源协议

Apache 2.0 License
