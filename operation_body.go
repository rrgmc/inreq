package inreq

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"reflect"
)

// DecodeOperationBody is a DecodeOperation that reads values from the request body.
type DecodeOperationBody struct {
}

func (d *DecodeOperationBody) Decode(ctx DecodeContext, r *http.Request, field reflect.Value,
	typ reflect.Type, tag *Tag) (bool, any, error) {
	if ctx.IsBodyDecoded() {
		return false, nil, fmt.Errorf("body was already decoded for type '%s'", typ.String())
	}
	fv := field
	if fv.CanAddr() {
		fv = fv.Addr()
	}
	found, err := decodeBody(ctx, r, fv.Interface(), tag)
	return found, IgnoreDecodeValue, err
}

func decodeBody(ctx DecodeContext, r *http.Request, data interface{}, tag *Tag) (bool, error) {
	if r.Body == nil {
		return false, nil
	}

	if !ctx.AllowReadBody() {
		return false, errors.New("body operation not allowed")
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return false, err
	}
	defer r.Body.Close()

	ctx.DecodedBody() // signal that the body was decoded

	if len(b) == 0 {
		return false, nil
	}

	var mediatype string

	if tag != nil {
		if typeStr := tag.Options.Value("type", ""); typeStr != "" {
			switch typeStr {
			case "json":
				mediatype = "application/json"
			case "xml":
				mediatype = "application/xml"
			default:
				return false, fmt.Errorf("invalid body type: '%s'", typeStr)
			}
		}
	}

	if mediatype == "" {
		mediatype, _, err = mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			return false, fmt.Errorf("error detecting body content type: %w", err)
		}
	}

	switch mediatype {
	case "application/json":
		err := json.Unmarshal(b, &data)
		if err != nil {
			return true, fmt.Errorf("error parsing JSON body: %w", err)
		}
		return true, nil
	case "text/xml", "application/xml":
		err := xml.Unmarshal(b, &data)
		if err != nil {
			return true, fmt.Errorf("error parsing XML body: %w", err)
		}
		return true, nil
	}

	return false, nil
}
