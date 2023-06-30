package inreq

import (
	"net/http"

	"github.com/RangelReale/instruct"
)

const (
	OperationIgnore  string = instruct.OperationIgnore
	OperationRecurse        = instruct.OperationRecurse
)

// Default operations.
const (
	OperationQuery  string = "query"
	OperationPath          = "path"
	OperationHeader        = "header"
	OperationForm          = "form"
	OperationBody          = "body"
)

// DecodeOperation is the interface for the http request-to-struct decoders.
type DecodeOperation = instruct.DecodeOperation[*http.Request, DecodeContext]

// IgnoreDecodeValue can be returned from [DecodeOperation.Decode] to signal that the value should not be set on the
// struct field. This is used for example in the "body" decoder.
var IgnoreDecodeValue = instruct.IgnoreDecodeValue
