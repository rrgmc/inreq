package inreq_test

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

func ExampleDecode() {
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
