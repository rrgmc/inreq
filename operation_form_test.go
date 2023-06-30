package inreq

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeForm(t *testing.T) {
	tests := []struct {
		name    string
		form    [][]string
		data    interface{}
		want    interface{}
		options []Option
		wantErr bool
	}{
		{
			name: "decode form",
			form: [][]string{{"val", "x1"}},
			data: &struct {
				Val string `inreq:"form"`
			}{},
			want: &struct {
				Val string `inreq:"form"`
			}{
				Val: "x1",
			},
		},
		{
			name: "decode form with slice",
			form: [][]string{{"val", "5", "6", "7"}},
			data: &struct {
				Val []int32 `inreq:"form"`
			}{},
			want: &struct {
				Val []int32 `inreq:"form"`
			}{
				Val: []int32{5, 6, 7},
			},
		},
		{
			name: "decode form with name",
			form: [][]string{{"XVal", "x1"}},
			data: &struct {
				Val string `inreq:"form,name=XVal"`
			}{},
			want: &struct {
				Val string `inreq:"form,name=XVal"`
			}{
				Val: "x1",
			},
		},
		{
			name: "decode form with name error",
			form: [][]string{{"val", "x1"}},
			data: &struct {
				Val string `inreq:"form,name=XVal"`
			}{},
			wantErr: true,
		},
		{
			name: "decode form not used values error",
			form: [][]string{{"val", "x1"}, {"val2", "x2"}},
			data: &struct {
				Val string `inreq:"form"`
			}{},
			wantErr: true,
			options: []Option{
				WithEnsureAllFormUsed(true),
			},
		},
		{
			name: "decode form recursive",
			form: [][]string{{"val", "x1"}},
			data: &struct {
				Inner struct {
					Val string `inreq:"form"`
				} `inreq:"recurse"`
			}{},
			want: &struct {
				Inner struct {
					Val string `inreq:"form"`
				} `inreq:"recurse"`
			}{
				Inner: struct {
					Val string `inreq:"form"`
				}{
					Val: "x1",
				},
			},
		},
		{
			name: "decode form with map tags",
			form: [][]string{{"val", "x1"}},
			data: &struct {
				Val string
			}{},
			want: &struct {
				Val string
			}{
				Val: "x1",
			},
			options: []Option{
				WithMapTags(map[string]any{
					"Val": "form",
				}),
			},
		},
		{
			name: "decode form with unused map tags error",
			form: [][]string{{"val", "x1"}},
			data: &struct {
				Val string
			}{},
			wantErr: true,
			options: []Option{
				WithMapTags(map[string]any{
					"Val":     "form",
					"NOTUSED": "nothing",
				}),
			},
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", nil)
			err := r.ParseForm()
			require.NoError(t, err)

			for _, qvalue := range tt.form {
				for _, v := range qvalue[1:] {
					r.Form.Add(qvalue[0], v)
				}
			}

			options := append(append([]Option{}, tt.options...),
				WithDecodeOperation(OperationForm, &DecodeOperationForm{}),
			)

			err = CustomDecode(r, tt.data, options...)
			if !tt.wantErr {
				require.NoError(t, err)
				require.Equal(t, tt.want, tt.data)
			} else if err == nil {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	require.NoError(t, nil)
}
