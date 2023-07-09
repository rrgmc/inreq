package inreq

import (
	"fmt"
	"net/http"
	"reflect"
)

// DecodeOperationPath is a DecodeOperation that gets values from HTTP paths (or routes).
// This is always framework-specific.
type DecodeOperationPath struct {
}

func (d *DecodeOperationPath) Decode(ctx DecodeContext, r *http.Request, isList bool, field reflect.Value,
	tag *Tag) (bool, any, error) {
	if ctx.PathValue() == nil {
		return false, nil, fmt.Errorf("path value function not set")
	}

	ctx.ValueUsed(OperationPath, tag.Name)
	return ctx.PathValue().GetRequestPath(r, tag.Name)
}
