package inreq

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeHeader(t *testing.T) {
	tests := []struct {
		name    string
		headers [][]string
		data    interface{}
		want    interface{}
		options []Option
		wantErr bool
	}{
		{
			name:    "decode header",
			headers: [][]string{{"Val", "x1"}},
			data: &struct {
				Val string `inreq:"header"`
			}{},
			want: &struct {
				Val string `inreq:"header"`
			}{
				Val: "x1",
			},
		},
		{
			name:    "decode header with slice",
			headers: [][]string{{"Val", "5", "6", "7"}},
			data: &struct {
				Val []int32 `inreq:"header"`
			}{},
			want: &struct {
				Val []int32 `inreq:"header"`
			}{
				Val: []int32{5, 6, 7},
			},
		},
		{
			name:    "decode header with name",
			headers: [][]string{{"XVal", "x1"}},
			data: &struct {
				Val string `inreq:"header,name=XVal"`
			}{},
			want: &struct {
				Val string `inreq:"header,name=XVal"`
			}{
				Val: "x1",
			},
		},
		{
			name:    "decode header with name error",
			headers: [][]string{{"Val", "x1"}},
			data: &struct {
				Val string `inreq:"header,name=XVal"`
			}{},
			wantErr: true,
		},
		{
			name:    "decode header recursive",
			headers: [][]string{{"Val", "x1"}},
			data: &struct {
				Inner struct {
					Val string `inreq:"header"`
				} `inreq:"recurse"`
			}{},
			want: &struct {
				Inner struct {
					Val string `inreq:"header"`
				} `inreq:"recurse"`
			}{
				Inner: struct {
					Val string `inreq:"header"`
				}{
					Val: "x1",
				},
			},
		},
		{
			name:    "decode header with map tags",
			headers: [][]string{{"Val", "x1"}},
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
					"Val": "header",
				}),
			},
		},
		{
			name:    "decode header with unused map tags error",
			headers: [][]string{{"Val", "x1"}},
			data: &struct {
				Val string
			}{},
			wantErr: true,
			options: []Option{
				WithMapTags(map[string]any{
					"Val":     "header",
					"NOTUSED": "nothing",
				}),
			},
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", nil)
			for _, qvalue := range tt.headers {
				for _, v := range qvalue[1:] {
					r.Header.Add(qvalue[0], v)
				}
			}

			options := append(append([]Option{}, tt.options...),
				WithDecodeOperation(OperationHeader, &DecodeOperationHeader{}),
			)

			err := CustomDecode(r, tt.data, options...)
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
