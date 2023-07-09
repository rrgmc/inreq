package inreq

import (
	"net/http"
	"reflect"
	"strings"

	"golang.org/x/exp/maps"
)

// DecodeOperationQuery is a DecodeOperation that gets values from HTTP query parameters.
type DecodeOperationQuery struct {
}

func (d *DecodeOperationQuery) Decode(ctx DecodeContext, r *http.Request, isList bool, field reflect.Value,
	tag *Tag) (bool, any, error) {
	if !r.URL.Query().Has(tag.Name) {
		return false, nil, nil
	}

	if isList {
		explode, err := tag.Options.BoolValue("explode", false)
		if err != nil {
			return false, nil, err
		}

		var value []string
		if explode {
			value = strings.Split(r.URL.Query().Get(tag.Name),
				tag.Options.Value("explodesep", ctx.SliceSplitSeparator()))
		} else {
			value = r.URL.Query()[tag.Name]
		}

		ctx.ValueUsed(OperationQuery, tag.Name)
		return true, value, nil
	}

	ctx.ValueUsed(OperationQuery, tag.Name)
	return true, r.URL.Query().Get(tag.Name), nil
}

func (d *DecodeOperationQuery) Validate(ctx DecodeContext, r *http.Request) error {
	if !ctx.EnsureAllQueryUsed() {
		return nil
	}

	queryKeys := map[string]bool{}
	for key, _ := range r.URL.Query() {
		queryKeys[key] = true
	}

	if !maps.Equal(queryKeys, ctx.GetUsedValues(OperationQuery)) {
		return ValuesNotUsedError{Operation: OperationQuery}
	}

	return nil
}
