package inreq

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodePath(t *testing.T) {
	tests := []struct {
		name       string
		pathValues [][2]string
		data       interface{}
		want       interface{}
		options    []AnyOption
		wantErr    bool
	}{
		{
			name:       "decode path",
			pathValues: [][2]string{{"val", "x1"}},
			data: &struct {
				Val string `inreq:"path"`
			}{},
			want: &struct {
				Val string `inreq:"path"`
			}{
				Val: "x1",
			},
		},
		{
			name:       "decode path with slice",
			pathValues: [][2]string{{"val", "5,6,7"}},
			data: &struct {
				Val []int32 `inreq:"path"`
			}{},
			wantErr: true,
		},
		{
			name:       "decode path with name",
			pathValues: [][2]string{{"XVal", "x1"}},
			data: &struct {
				Val string `inreq:"path,name=XVal"`
			}{},
			want: &struct {
				Val string `inreq:"path,name=XVal"`
			}{
				Val: "x1",
			},
		},
		{
			name:       "decode path with name error",
			pathValues: [][2]string{{"val", "x1"}},
			data: &struct {
				Val string `inreq:"path,name=XVal"`
			}{},
			wantErr: true,
		},
		{
			name:       "decode path recursive",
			pathValues: [][2]string{{"val", "x1"}},
			data: &struct {
				Inner struct {
					Val string `inreq:"path"`
				} `inreq:"recurse"`
			}{},
			want: &struct {
				Inner struct {
					Val string `inreq:"path"`
				} `inreq:"recurse"`
			}{
				Inner: struct {
					Val string `inreq:"path"`
				}{
					Val: "x1",
				},
			},
		},
		{
			name:       "decode path with map tags",
			pathValues: [][2]string{{"val", "x1"}},
			data: &struct {
				Val string
			}{},
			want: &struct {
				Val string
			}{
				Val: "x1",
			},
			options: []AnyOption{
				WithMapTags(map[string]any{
					"Val": "path",
				}),
			},
		},
		{
			name:       "decode path with unused map tags error",
			pathValues: [][2]string{{"val", "x1"}},
			data: &struct {
				Val string
			}{},
			wantErr: true,
			options: []AnyOption{
				WithMapTags(map[string]any{
					"Val":     "path",
					"NOTUSED": "nothing",
				}),
			},
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", nil)

			options := append(append([]AnyOption{}, tt.options...),
				WithDecodeOperation(OperationPath, &DecodeOperationPath{}),
				WithPathValue(PathValueFunc(func(r *http.Request, name string) (found bool, value any, err error) {
					for _, qvalue := range tt.pathValues {
						if name == qvalue[0] {
							return true, qvalue[1], nil
						}
					}
					return false, nil, nil
				})),
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
