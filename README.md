# go-bulkrequest 

---

## Description

**go-bulkrequest** is a simple library for fetching and parsing data from multiple URLs/endpoints.

Project is developed in [go](https://golang.org/).

Every bulk request is highly configurable, for example:
- You can use proxy connection instead of direct connection, for every proxy, requests will be parallelized.
- You can set cookies/headers/user-agent for every request.
- You can provide custom parser for every request.

## Installation

```bash
go get github.com/jodua/go-bulkrequest
```

## Usage

### Importing library

```go
import "github.com/jodua/go-bulkrequest/bulkrequest"
```

### Creating bulk request builder

```go
bulkRequest := bulkrequest.NewBulkRequest()
```

### Customizing bulk request

**Builder methods**

SetBaseUrl - sets base URL for all requests.

_Params:_
- `baseUrl string` - baseUrl

```go
bulkRequest.SetBaseUrl("https://example.com")
```

SetUrls - sets list of URL suffixes for all requests.

_Params:_
- `urls []string` - list of URL suffixes

Following example will create requests for URLs:
- Base URL + Urls[0]
- Base URL + Urls[1]
- ...
- Base URL + Urls[len(Urls)-1]

```go
urlList := []string{"1", "2", "3", "4", "5"}
bulkRequest.SetUrls(urlList)
```

SetTimeout - sets timeout after which request will be aborted.

_Params:_
- `timeout time.Duration` - duration of timeout

```go
bulkRequest.SetTimeout(time.Second * 10)
````

AddHeader - adds header to all requests.

_Params:_
- `key string` - header key
- `value string` - header value

```go
bulkRequest.AddHeader("X-Header", "Value")
```

AddCookie - adds cookie to all requests.

_Params:_
- `name string` - header key
- `value string` - header value

```go
bulkRequest.AddCookie("Cookie", "Value")
```

SetDelayConfig - sets delay configuration for all requests.

_Params:_
- `delayConfig *datatypes.DelayConfig` - delay configuration


Delay config is a struct that contains delay configuration for all requests.
It consists of:
- `DelayMin time.Duration` - minimum delay between requests
- `DelayMax time.Duration` - maximum delay between requests

When fetching data from multiple URLs, requests will be delayed between `DelayMin` and `DelayMax` time.

```go
import "github.com/jodua/go-bulkrequest/datatypes"

delayConfig := datatypes.DelayConfig{
    DelayMin: time.Second * 1,
    DelayMax: time.Second * 2,
}

bulkRequest.SetDelayConfig(&delayConfig)
```

SetProxyConfig - sets proxy configuration for all requests.
If proxy configuration is not set, requests will be made directly.

_Params:_
- `proxyConfig *datatypes.ProxyConfig` - proxy configuration


Proxy config is a struct that contains proxy configuration for all requests.
It consists of:
- `ProxyList []string` - list of proxies 
- `RequestsPerProxy` - number of requests that will be sent through each proxy

```go
import "github.com/jodua/go-bulkrequest/datatypes"

proxyList := []string{"http://proxy1:1231", "http://proxy2:1111"}

proxyConfig := datatypes.ProxyConfig{
    ProxyList:        proxyList,
    RequestsPerProxy: 5,
}

bulkRequest.SetProxyConfig(&delayConfig)
```

SetUserAgentConfig - sets user agent configuration for all requests.

_Params:_
- `userAgentConfig *datatypes.UserAgentConfig` - user agent configuration

User agent config is a struct that contains user agent configuration for all requests.
It consists of:
- `UserAgentList []string` - list of user agents
 
```go
import "github.com/jodua/go-bulkrequest/datatypes" 

userAgentList := []string{"UserAgent1", "UserAgent2"}

userAgentConfig := datatypes.UserAgentConfig{
    UserAgentList: userAgentList,
}

bulkRequest.SetUserAgentConfig(&userAgentConfig)
```

SetParser - sets parser for all requests.

_Params:_
- `parser *jsonparser.JSONParser` - pointer to parser object

Example parser can be found in `github.com/jodua/go-bulkrequest/JSONParser/schemas` package.

```go
var parser = schemas.JSONPlaceholderTodoParser

bulkRequest.SetParser(&parser)
```

### Building and executing bulk request

```go
bulkRequest := bulkRequest.Build()

fetch, err := bulkRequest.Fetch()
if err != nil {
    // handle error
}
log.Println(fetch)
```

Full example can be found in `main.go` file.

## Building JSON parser

JSONParser struct consists of:
- `JSONSchema any` - pointer to JSON schema object
- `ConvertFunction func(any,string) any` - function that converts JSON data to desired format, second parameter is request URL
- `ValidatorFunction func(any) error` - function that validates JSON data
- `Output any` - pointer to output struct that will be filled with data
- `Name string` - name of parser

Example parser can be found in `github.com/jodua/go-bulkrequest/JSONParser/schemas` package.

## Issues

File issues through **Issues** tab.

## License

MIT License

