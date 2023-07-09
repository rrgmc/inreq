package inreq

import (
	"mime/multipart"
	"net/http"
	"reflect"

	"golang.org/x/exp/maps"
)

// DecodeOperationForm is a DecodeOperation that gets values from HTTP forms.
type DecodeOperationForm struct {
}

func (d *DecodeOperationForm) Decode(ctx DecodeContext, r *http.Request, isList bool, field reflect.Value,
	tag *Tag) (bool, any, error) {
	var form multipart.Form

	err := r.ParseForm()
	if err != nil {
		return false, nil, err
	}

	if r.MultipartForm != nil {
		form = *r.MultipartForm
	} else {
		if r.Form != nil {
			form.Value = r.Form
		}
	}

	values, ok := form.Value[tag.Name]
	if !ok {
		return false, nil, nil
	}

	if len(values) == 0 {
		return false, nil, nil
	}

	ctx.ValueUsed(OperationForm, tag.Name)

	if isList {
		return true, values, nil
	}
	return true, values[0], nil
}

func (d *DecodeOperationForm) Validate(ctx DecodeContext, r *http.Request) error {
	if !ctx.EnsureAllFormUsed() {
		return nil
	}

	var form multipart.Form
	if r.MultipartForm != nil {
		form = *r.MultipartForm
	} else {
		if r.Form != nil {
			form.Value = r.Form
		}
	}

	formKeys := map[string]bool{}
	for key, _ := range form.Value {
		formKeys[key] = true
	}

	if !maps.Equal(formKeys, ctx.GetUsedValues(OperationForm)) {
		return ValuesNotUsedError{Operation: OperationForm}
	}

	return nil
}
