package inreq

import (
	"encoding"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"reflect"
)

var (
	textUnmarshalerType = reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()
)

// DecodeOperationBody is a DecodeOperation that reads values from the request body.
type DecodeOperationBody struct {
}

func (d *DecodeOperationBody) Decode(ctx DecodeContext, r *http.Request, isList bool, field reflect.Value,
	tag *Tag) (bool, any, error) {

	if r.Body == nil {
		return false, nil, nil
	}

	if !ctx.AllowReadBody() {
		return false, nil, errors.New("body operation not allowed")
	}
	if ctx.IsBodyDecoded() {
		return false, nil, fmt.Errorf("body was already decoded")
	}

	return decodeBody(ctx, r, field, tag)
}

func decodeBodyReadData(ctx DecodeContext, r *http.Request) (bool, []byte, error) {
	if ctx.IsBodyDecoded() {
		return true, nil, fmt.Errorf("body was already decoded")
	}

	ctx.DecodedBody() // signal that the body was decoded

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return true, nil, err
	}
	defer r.Body.Close()

	if len(b) == 0 {
		return false, nil, nil
	}

	return true, b, nil
}

func decodeBodyRaw(ctx DecodeContext, r *http.Request, field reflect.Value) (bool, bool, any, error) {
	// check for raw data
	if field.Type().PkgPath() == "" { // only if predeclared type (ignores for example "type IP []byte".)
		switch field.Type().Kind() {
		case reflect.String:
			rfound, rvalue, rerr := decodeBodyReadData(ctx, r)
			return true, rfound, string(rvalue), rerr
		case reflect.Slice:
			if field.Type().Elem().Kind() == reflect.Uint8 {
				rfound, rvalue, rerr := decodeBodyReadData(ctx, r)
				return true, rfound, rvalue, rerr
			}
		}
	}
	return false, false, nil, nil
}

func decodeBody(ctx DecodeContext, r *http.Request, field reflect.Value, tag *Tag) (bool, any, error) {
	// check for raw data
	rfound, found, data, err := decodeBodyRaw(ctx, r, field)
	if rfound {
		return found, data, err
	}

	// decode into struct field
	fv := field
	if fv.CanAddr() {
		fv = fv.Addr()
	}

	found, data, err = ctx.BodyDecoder().Unmarshal(ctx, tag.Options.Value("type", ""),
		r, fv.Interface())
	if found {
		return found, data, err
	}

	// try known interfaces
	if reflect.PointerTo(field.Type()).Implements(textUnmarshalerType) {
		// encoding.TextUnmarshaler
		rfound, rvalue, rerr := decodeBodyReadData(ctx, r)
		if rfound {
			if err != nil {
				return rfound, rvalue, rerr
			}

			xtarget := reflect.New(field.Type())
			um := xtarget.Interface().(encoding.TextUnmarshaler)
			if err := um.UnmarshalText(rvalue); err != nil {
				return true, nil, err
			}
			field.Set(xtarget.Elem())

			return true, IgnoreDecodeValue, nil
		}
	}

	return false, nil, nil
}

// defaultBodyDecoder decodes JSON and XML to structs.
type defaultBodyDecoder struct {
}

func NewDefaultBodyDecoder() BodyDecoder {
	return &defaultBodyDecoder{}
}

func (d defaultBodyDecoder) Unmarshal(ctx DecodeContext, typeParam string, r *http.Request, data any) (bool, any, error) {
	var mediaType string

	if typeParam != "" {
		switch typeParam {
		case "json":
			mediaType = "application/json"
		case "xml":
			mediaType = "application/xml"
		}
	}

	if mediaType == "" {
		var err error
		contentType := r.Header.Get("Content-Type")
		if contentType != "" {
			mediaType, _, err = mime.ParseMediaType(r.Header.Get("Content-Type"))
			if err != nil {
				return false, nil, fmt.Errorf("error detecting body content type: %w", err)
			}
		}
	}

	switch mediaType {
	case "application/json":
		ctx.DecodedBody()
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			return true, nil, fmt.Errorf("error parsing JSON body: %w", err)
		}
		return true, IgnoreDecodeValue, nil
	case "text/xml", "application/xml":
		ctx.DecodedBody()
		err := xml.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			return true, nil, fmt.Errorf("error parsing XML body: %w", err)
		}
		return true, IgnoreDecodeValue, nil
	}

	return false, nil, nil
}
