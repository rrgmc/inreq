package inreq

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeQuery(t *testing.T) {
	tests := []struct {
		name    string
		query   [][2]string
		data    interface{}
		want    interface{}
		options []Option
		wantErr bool
	}{
		{
			name:  "decode query",
			query: [][2]string{{"val", "x1"}},
			data: &struct {
				Val string `inreq:"query"`
			}{},
			want: &struct {
				Val string `inreq:"query"`
			}{
				Val: "x1",
			},
		},
		{
			name:  "decode query with slice",
			query: [][2]string{{"val", "5,6,7"}},
			data: &struct {
				Val []int32 `inreq:"query"`
			}{},
			want: &struct {
				Val []int32 `inreq:"query"`
			}{
				Val: []int32{5, 6, 7},
			},
		},
		{
			name:  "decode query with name",
			query: [][2]string{{"XVal", "x1"}},
			data: &struct {
				Val string `inreq:"query,name=XVal"`
			}{},
			want: &struct {
				Val string `inreq:"query,name=XVal"`
			}{
				Val: "x1",
			},
		},
		{
			name:  "decode query with name error",
			query: [][2]string{{"val", "x1"}},
			data: &struct {
				Val string `inreq:"query,name=XVal"`
			}{},
			wantErr: true,
		},
		{
			name:  "decode query not used values error",
			query: [][2]string{{"val", "x1"}, {"val2", "x2"}},
			data: &struct {
				Val string `inreq:"query"`
			}{},
			wantErr: true,
			options: []Option{
				WithEnsureAllQueryUsed(true),
			},
		},
		{
			name:  "decode query recursive",
			query: [][2]string{{"val", "x1"}},
			data: &struct {
				Inner struct {
					Val string `inreq:"query"`
				} `inreq:"recurse"`
			}{},
			want: &struct {
				Inner struct {
					Val string `inreq:"query"`
				} `inreq:"recurse"`
			}{
				Inner: struct {
					Val string `inreq:"query"`
				}{
					Val: "x1",
				},
			},
		},
		{
			name:  "decode query with map tags",
			query: [][2]string{{"val", "x1"}},
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
					"Val": "query",
				}),
			},
		},
		{
			name:  "decode query with unused map tags error",
			query: [][2]string{{"val", "x1"}},
			data: &struct {
				Val string
			}{},
			wantErr: true,
			options: []Option{
				WithMapTags(map[string]any{
					"Val":     "query",
					"NOTUSED": "nothing",
				}),
			},
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", nil)
			q := r.URL.Query()
			for _, qvalue := range tt.query {
				q.Add(qvalue[0], qvalue[1])
			}
			r.URL.RawQuery = q.Encode()

			options := append(append([]Option{}, tt.options...),
				WithDecodeOperation(OperationQuery, &DecodeOperationQuery{}),
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
