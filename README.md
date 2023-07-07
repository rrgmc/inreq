# InReq - Golang http request to struct
[![GoDoc](https://godoc.org/github.com/RangelReale/inreq?status.png)](https://godoc.org/github.com/RangelReale/inreq)

InReq is a Golang library to extract information from `*http.Request` into structs. It does this using 
struct tags and/or a configuration map.

It is highly configurable: configurations can be entirely in maps without requiring struct changes, custom decoders 
can be created, configurations can be overriden on specific calls, a field name mapper can be set, custom type 
resolvers are available (or the entire type resolving logic can be replaced), the HTTP body can be parsed into 
a specific field, and much more.

## Examples

```go
import (
    "fmt"
    "net/http"
    "strings"

    "github.com/RangelReale/inreq"
)

type InputBody struct {
    DeviceID string `json:"device_id"`
    Name     string `json:"name"`
}

type Input struct {
    AuthToken      string    `inreq:"header,name=X-Auth-Token"`
    DeviceID       string    `inreq:"path"`
    WithDetails    bool      `inreq:"query,name=with_details"`
    Page           int       `inreq:"query"`
    Body           InputBody `inreq:"body"`
    FormDeviceName string    `inreq:"form,name=devicename"`
}

func main() {
    r, err := http.NewRequest(http.MethodPost, "/device/12345?with_details=true&page=2",
        strings.NewReader(`{"device_id":"12345","name":"Device for testing"}`))
    if err != nil {
        panic(err)
    }
    err = r.ParseForm()
    if err != nil {
        panic(err)
    }
    r.Header.Add("Content-Type", "application/json")
    r.Header.Add("X-Auth-Token", "auth-token-value")
    r.Form.Add("devicename", "form-device-name")

    data := &Input{}

    err = inreq.Decode(r, data,
        // usually this will be a framework-specific implementation, like "github.com/RangelReale/inreq-path/gorillamux".
        inreq.WithPathValue(inreq.PathValueFunc(func(r *http.Request, name string) (found bool, value any, err error) {
            if name == "deviceid" {
                return true, "12345", err
            }
            return false, nil, nil
        })))
    if err != nil {
        panic(err)
    }

    fmt.Printf("Auth Token: %s\n", data.AuthToken)
    fmt.Printf("Device ID: %s\n", data.DeviceID)
    fmt.Printf("With details: %t\n", data.WithDetails)
    fmt.Printf("Page: %d\n", data.Page)
    fmt.Printf("Body Device ID: %s\n", data.Body.DeviceID)
    fmt.Printf("Body Name: %s\n", data.Body.Name)
    fmt.Printf("Form Device Name: %s\n", data.FormDeviceName)

    // Output: Auth Token: auth-token-value
    // Device ID: 12345
    // With details: true
    // Page: 2
    // Body Device ID: 12345
    // Body Name: Device for testing
    // Form Device Name: form-device-name
}
```

Using generics:

```go
import (
    "fmt"
    "net/http"
    "strings"

    "github.com/RangelReale/inreq"
)

type InputTypeBody struct {
    DeviceID string `json:"device_id"`
    Name     string `json:"name"`
}

type InputType struct {
    AuthToken      string        `inreq:"header,name=X-Auth-Token"`
    DeviceID       string        `inreq:"path"`
    WithDetails    bool          `inreq:"query,name=with_details"`
    Page           int           `inreq:"query"`
    Body           InputTypeBody `inreq:"body"`
    FormDeviceName string        `inreq:"form,name=devicename"`
}

func main() {
    r, err := http.NewRequest(http.MethodPost, "/device/12345?with_details=true&page=2",
        strings.NewReader(`{"device_id":"12345","name":"Device for testing"}`))
    if err != nil {
        panic(err)
    }
    err = r.ParseForm()
    if err != nil {
        panic(err)
    }
    r.Header.Add("Content-Type", "application/json")
    r.Header.Add("X-Auth-Token", "auth-token-value")
    r.Form.Add("devicename", "form-device-name")

    data, err := inreq.DecodeType[InputType](r,
        // usually this will be a framework-specific implementation, like "github.com/RangelReale/inreq-path/gorillamux".
        inreq.WithPathValue(inreq.PathValueFunc(func(r *http.Request, name string) (found bool, value any, err error) {
            if name == "deviceid" {
                return true, "12345", err
            }
            return false, nil, nil
        })))
    if err != nil {
        panic(err)
    }

    fmt.Printf("Auth Token: %s\n", data.AuthToken)
    fmt.Printf("Device ID: %s\n", data.DeviceID)
    fmt.Printf("With details: %t\n", data.WithDetails)
    fmt.Printf("Page: %d\n", data.Page)
    fmt.Printf("Body Device ID: %s\n", data.Body.DeviceID)
    fmt.Printf("Body Name: %s\n", data.Body.Name)
    fmt.Printf("Form Device Name: %s\n", data.FormDeviceName)

    // Output: Auth Token: auth-token-value
    // Device ID: 12345
    // With details: true
    // Page: 2
    // Body Device ID: 12345
    // Body Name: Device for testing
    // Form Device Name: form-device-name
}
```

## Default operations

### query

`inreq:"query,name=<query-param-name>,required=true,explode=false,explodesep=,"`

- name: the query parameter name to get from `req.URL.Query().Get()`. Default uses `FieldNameMapper`, which by default uses `strings.ToLower`.
- required: whether the query parameter is required to exist. Default is true.
- explode: whether to use `strings.Split` on the query string if the target struct field is a slice. Default is false.
- explodesep: the separator to use when exploding the string.

### header

`inreq:"header,name=<header-name>,required=true"`

- name: the header name to get from `req.Header.Values()`. Default uses `FieldNameMapper`, which by default uses `strings.ToLower`.
- required: whether the header is required to exist. Default is true.

### form

`inreq:"form,name=<form-field-name>,required=true"`

- name: the form field name to get from `req.Form.Get()` or `req.MultipartForm.Value`. Default uses `FieldNameMapper`, which by default uses `strings.ToLower`.
- required: whether the form field is required to exist. Default is true.

### path

`inreq:"path,name=<path-var-name>,required=true"`

A path isn't an HTTP concept, but usually http frameworks have a concept of `routes` which can contain path variables,
a framework-specific function should be set using `WithPathValue`. Some of these are available a
[https://github.com/RangelReale/inreq-path](https://github.com/RangelReale/inreq-path).

- name: the path var name to get from `PathValue.GetRequestPath`. Default uses `FieldNameMapper`, which by default uses `strings.ToLower`.
- required: whether the path var is required to exist. Default is true.

### body

`inreq:"body,required=true,type=json"`

Body unmarshals data into the struct field, usually JSON or XML.

- required: whether an HTTP body required to exist. Default is true.
- type: type of body to decode. If blank, will use the `Content-Type` header. Should be only a type name ("json", "xml").

# Author

Rangel Reale (rangelreale@gmail.com)
