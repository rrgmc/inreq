package inreq

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeBody(t *testing.T) {
	tests := []struct {
		name             string
		headers          [][]string
		skipContentType  bool
		body             string
		data             interface{}
		want             interface{}
		options          []Option
		wantErr          bool
		wantCompareError bool
	}{
		{
			name: "decode root body",
			body: `{"Val": "x1"}`,
			data: &struct {
				_   StructOption `inreq:"body"`
				Val string
			}{},
			want: &struct {
				_   StructOption `inreq:"body"`
				Val string
			}{
				Val: "x1",
			},
		},
		{
			name:    "decode root body invalid operation",
			headers: [][]string{{"S", "x2"}},
			body:    `{"Val": "x1"}`,
			data: &struct {
				_   StructOption `inreq:"header,name=S"`
				Val string
			}{},
			wantErr: true,
		},
		{
			name: "decode root body compare error",
			body: `{"Val": "x2"}`,
			data: &struct {
				_   StructOption `inreq:"body"`
				Val string
			}{},
			want: &struct {
				_   StructOption `inreq:"body"`
				Val string
			}{
				Val: "x1",
			},
			wantCompareError: true,
		},
		{
			name:    "decode root body before",
			headers: [][]string{{"Val", "x2"}},
			body:    `{"Val": "x1"}`,
			data: &struct {
				_   StructOption `inreq:"body,so_recurse=true,so_when=before"`
				Val string       `inreq:"header"`
			}{},
			want: &struct {
				_   StructOption `inreq:"body,so_recurse=true,so_when=before"`
				Val string       `inreq:"header"`
			}{
				Val: "x2",
			},
		},
		{
			name:    "decode root body after",
			headers: [][]string{{"Val", "x2"}},
			body:    `{"Val": "x1"}`,
			data: &struct {
				_   StructOption `inreq:"body,so_recurse=true,so_when=after"`
				Val string       `inreq:"header"`
			}{},
			want: &struct {
				_   StructOption `inreq:"body,so_recurse=true,so_when=after"`
				Val string       `inreq:"header"`
			}{
				Val: "x1",
			},
		},
		{
			name: "decode root body using map tags",
			body: `{"Val": "x1"}`,
			data: &struct {
				_   StructOption
				Val string
			}{},
			want: &struct {
				_   StructOption
				Val string
			}{
				Val: "x1",
			},
			options: []Option{
				WithMapTags(map[string]any{
					StructOptionMapTag: "body",
				}),
			},
		},
		{
			name: "decode body field",
			body: `{"Val": "x1"}`,
			data: &struct {
				B struct {
					Val string
				} `inreq:"body"`
			}{},
			want: &struct {
				B struct {
					Val string
				} `inreq:"body"`
			}{
				B: struct {
					Val string
				}{
					Val: "x1",
				},
			},
		},
		{
			name: "decode body field not allowed error",
			body: `{"Val": "x1"}`,
			data: &struct {
				B struct {
					Val string
				} `inreq:"body"`
			}{},
			wantErr: true,
			options: []Option{
				WithAllowReadBody(false),
			},
		},
		{
			name: "decode body field compare error",
			body: `{"Val": "x2"}`,
			data: &struct {
				B struct {
					Val string
				} `inreq:"body"`
			}{},
			want: &struct {
				B struct {
					Val string
				} `inreq:"body"`
			}{
				B: struct {
					Val string
				}{
					Val: "x1",
				},
			},
			wantCompareError: true,
		},
		{
			name: "decode body field multiple error",
			body: `{"Val": "x1"}`,
			data: &struct {
				B struct {
					Val string
				} `inreq:"body"`
				C struct {
					Val string
				} `inreq:"body"`
			}{},
			wantErr: true,
		},
		{
			name:            "decode body field with type",
			skipContentType: true,
			body:            `{"Val": "x1"}`,
			data: &struct {
				B struct {
					Val string
				} `inreq:"body,type=json"`
			}{},
			want: &struct {
				B struct {
					Val string
				} `inreq:"body,type=json"`
			}{
				B: struct {
					Val string
				}{
					Val: "x1",
				},
			},
		},
		{
			name:            "decode body field without type error",
			skipContentType: true,
			body:            `{"Val": "x1"}`,
			data: &struct {
				B struct {
					Val string
				} `inreq:"body"`
			}{},
			wantErr: true,
		},
		{
			name: "decode body field invalid JSON",
			body: `A{@`,
			data: &struct {
				B struct {
					Val string
				} `inreq:"body"`
			}{},
			wantErr: true,
		},
		{
			name: "decode body with slice",
			body: `{"Val": [5,6,7]}`,
			data: &struct {
				B struct {
					Val []int32
				} `inreq:"body"`
			}{},
			want: &struct {
				B struct {
					Val []int32
				} `inreq:"body"`
			}{
				B: struct {
					Val []int32
				}{
					Val: []int32{5, 6, 7},
				},
			},
		},
		{
			name: "decode body with slice compare error",
			body: `{"Val": [5,6,8]}`,
			data: &struct {
				B struct {
					Val []int32
				} `inreq:"body"`
			}{},
			want: &struct {
				B struct {
					Val []int32
				} `inreq:"body"`
			}{
				B: struct {
					Val []int32
				}{
					Val: []int32{5, 6, 7},
				},
			},
			wantCompareError: true,
		},
		{
			name: "decode body recursive",
			body: `{"Val": "x1"}`,
			data: &struct {
				Inner struct {
					B struct {
						Val string
					} `inreq:"body"`
				} `inreq:"recurse"`
			}{},
			want: &struct {
				Inner struct {
					B struct {
						Val string
					} `inreq:"body"`
				} `inreq:"recurse"`
			}{
				Inner: struct {
					B struct {
						Val string
					} `inreq:"body"`
				}{
					B: struct {
						Val string
					}{
						Val: "x1",
					},
				},
			},
		},
		{
			name: "decode query with map tags",
			body: `{"Val": "x1"}`,
			data: &struct {
				B struct {
					Val string
				}
			}{},
			want: &struct {
				B struct {
					Val string
				}
			}{
				B: struct {
					Val string
				}{
					Val: "x1",
				},
			},
			options: []Option{
				WithMapTags(map[string]any{
					"B": "body",
				}),
			},
		},
		{
			name: "decode query with unused map tags error",
			body: `{"Val": "x1"}`,
			data: &struct {
				B struct {
					Val string
				}
			}{},
			want: &struct {
				B struct {
					Val string
				}
			}{
				B: struct {
					Val string
				}{
					Val: "x1",
				},
			},
			wantErr: true,
			options: []Option{
				WithMapTags(map[string]any{
					"B":       "body",
					"NOTUSED": "nothing",
				}),
			},
		},
		{
			name: "decode body field with inner struct option",
			body: `{"Val": "x1"}`,
			data: &struct {
				B struct {
					_   StructOption `inreq:"body"`
					Val string
				}
			}{},
			want: &struct {
				B struct {
					_   StructOption `inreq:"body"`
					Val string
				}
			}{
				B: struct {
					_   StructOption `inreq:"body"`
					Val string
				}{
					Val: "x1",
				},
			},
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			for _, qvalue := range tt.headers {
				for _, v := range qvalue[1:] {
					r.Header.Add(qvalue[0], v)
				}
			}
			if !tt.skipContentType {
				r.Header.Add("Content-Type", "application/json")
			}

			options := append(append([]Option{}, tt.options...),
				WithDecodeOperation(OperationBody, &DecodeOperationBody{}),
			)
			if len(tt.headers) > 0 {
				options = append(options,
					WithDecodeOperation(OperationHeader, &DecodeOperationHeader{}),
				)
			}

			err := CustomDecode(r, tt.data, options...)
			if !tt.wantErr {
				require.NoError(t, err)
				if !tt.wantCompareError {
					require.Equal(t, tt.want, tt.data)
				} else {
					require.NotEqual(t, tt.want, tt.data)
				}
			} else if err == nil {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	require.NoError(t, nil)
}
