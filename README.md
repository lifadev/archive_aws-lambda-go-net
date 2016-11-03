[<img src="_asset/powered-by-aws.png" alt="Powered by Amazon Web Services" align="right">][aws-home]
[<img src="_asset/created-by-eawsy.png" alt="Created by eawsy" align="right">][eawsy-home]

# eawsy/aws-lambda-go-net
> A seamless way to execute Go web applications on AWS Lambda and AWS API Gateway.

[![Runtime][runtime-badge]][eawsy-runtime]
[![Api][api-badge]][eawsy-godoc]
[![Chat][chat-badge]][eawsy-gitter]
![Status][status-badge]
[![License][license-badge]](LICENSE)
<sup>•</sup> <sup>•</sup> <sup>•</sup>
[![Hire us][hire-badge]][eawsy-hire-form]

[AWS Lambda][aws-lambda-home] lets you run code without provisioning or managing servers. 
[AWS API Gateway][aws-gtw-home] is a fully managed service that makes it easy for developers to create, publish, 
maintain, monitor, and secure APIs at any scale. 

This projects provides an AWS Lambda [network interface for Go][go-net-listener]. You can leverage AWS Lambda and 
AWS API Gateway to handle web requests using *any* Go application framework.

## Preview

> Below are some examples with popular Go frameworks. These examples are not exhaustive.
  **You are free to use anything you want!**

```sh
go get -u -d github.com/eawsy/aws-lambda-go-net/...
```

### Using [Go][go-http-pkg]

```go
package main

import (
  "net/http"

  "github.com/eawsy/aws-lambda-go-net/service/lambda/runtime/net"
)

func handle(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Hello, World!"))
}

func init() {
  go http.Serve(net.Listener(), http.HandlerFunc(handle))
}

func main() {}
```

### Using [Gin][gin-github]

```sh
go get gopkg.in/gin-gonic/gin.v1
```
```go
package main

import (
  "net/http"

  "github.com/eawsy/aws-lambda-go-net/service/lambda/runtime/net"

  "github.com/gin-gonic/gin"
)

func handle(ctx *gin.Context) {
  ctx.String(http.StatusOK, "Hello, %s!", ctx.Param("name"))
}

func init() {
  r := gin.Default()
  r.GET("/hello/:name", handle)
  go http.Serve(net.Listener(), r)
}

func main() {}
```

### Using [Iris][iris-github]

```sh
go get -u github.com/kataras/iris/iris
```
```go
package main

import (
  "github.com/eawsy/aws-lambda-go-net/service/lambda/runtime/net"
  "github.com/kataras/iris"
)

func handle(ctx *iris.Context) {
  ctx.Write("Hello, %s!", ctx.Param("name"))
}

func init() {
  iris.Get("/hello/:name", handle)
  go iris.Serve(net.Listener())
}

func main() {}
```

### Using [Gorilla][gorilla-github]

```sh
go get -u github.com/gorilla/mux
```
```go
package main

import (
  "fmt"
  "net/http"

  "github.com/eawsy/aws-lambda-go-net/service/lambda/runtime/net"

  "github.com/gorilla/mux"
)

func handle(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte(fmt.Sprintf("Hello, %s!", mux.Vars(r)["name"])))
}

func init() {
  r := mux.NewRouter()
  r.HandleFunc("/hello/{name}", handle)
  go http.Serve(net.Listener(), r)
}

func main() {}
```

## Documentation

This [wiki][eawsy-wiki] is the main source of documentation for developers working with or contributing to the 
project.

## About

[![eawsy](_asset/eawsy-logo.png)][eawsy-home]

This project is maintained and funded by Alsanium, SAS.

[We][eawsy-home] :heart: [AWS][aws-home] and open source software. See [our other projects][eawsy-github], or 
[hire us][eawsy-hire-form] to help you build modern applications on AWS.

## License

This product is licensed to you under the Apache License, Version 2.0 (the "License"); you may not use this product 
except in compliance with the License. See [LICENSE](LICENSE) and [NOTICE](NOTICE) for more information.

## Trademark

Alsanium, eawsy, the "Created by eawsy" logo, and the "eawsy" logo are trademarks of Alsanium, SAS. or its affiliates 
in France and/or other countries.

Amazon Web Services, the "Powered by Amazon Web Services" logo, and AWS Lambda are trademarks of Amazon.com, Inc. or 
its affiliates in the United States and/or other countries.

[eawsy-home]: https://eawsy.com
[eawsy-github]: https://github.com/eawsy
[eawsy-runtime]: https://github.com/eawsy/aws-lambda-go
[eawsy-gitter]: https://gitter.im/eawsy/bavardage
[eawsy-godoc]: https://godoc.org/github.com/eawsy/aws-lambda-go-net/service/lambda/runtime/net
[eawsy-wiki]: https://github.com/eawsy/aws-lambda-go-net/wiki
[eawsy-hire-form]: https://docs.google.com/forms/d/e/1FAIpQLSfPvn1Dgp95DXfvr3ClPHCNF5abi4D1grveT5btVyBHUk0nXw/viewform
[aws-home]: https://aws.amazon.com/
[aws-lambda-home]: https://aws.amazon.com/lambda/
[aws-gtw-home]: https://aws.amazon.com/api-gateway/
[go-net-listener]: https://golang.org/pkg/net/#Listener
[go-http-pkg]: https://golang.org/pkg/net/http/
[gin-github]: https://github.com/gin-gonic/gin
[iris-github]: https://github.com/kataras/iris
[gorilla-github]: https://github.com/gorilla/mux
[runtime-badge]: http://img.shields.io/badge/runtime-go-ef6c00.svg?style=flat-square
[api-badge]: http://img.shields.io/badge/api-godoc-7986cb.svg?style=flat-square
[chat-badge]: http://img.shields.io/badge/chat-gitter-e91e63.svg?style=flat-square
[status-badge]: http://img.shields.io/badge/status-beta-827717.svg?style=flat-square
[license-badge]: http://img.shields.io/badge/license-apache-757575.svg?style=flat-square
[hire-badge]: http://img.shields.io/badge/hire-eawsy-2196f3.svg?style=flat-square
