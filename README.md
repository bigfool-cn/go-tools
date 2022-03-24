## go-tools

Encapsulate golang tools operations

- http

#### Use
```shell
go get -u github.com/bigfool-cn/go-tools
```

#### example
> http
```go
package main
import "github.com/bigfool-cn/go-tools"

func main()  {
    httpClient := tools_http.NewHttpClient()
    statusCode, response, err := httpClient.SetClient(xxx).SetMethod(xxx).SetUrl(xxx).SetHeader(xxx).SetBody(xxx).Do()
}
```