package inreq

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeType(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	q := r.URL.Query()
	q.Add("val", "v1")
	r.URL.RawQuery = q.Encode()

	type DataType struct {
		Val string `inreq:"query"`
	}

	v, err := DecodeType[DataType](r)
	require.NoError(t, err)
	require.Equal(t, "v1", v.Val)
}

func TestCustomDecodeType(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	q := r.URL.Query()
	q.Add("val", "v1")
	r.URL.RawQuery = q.Encode()

	type DataType struct {
		Val string `inreq:"query"`
	}

	v, err := CustomDecodeType[DataType](r, WithDefaultDecodeOperations())
	require.NoError(t, err)
	require.Equal(t, "v1", v.Val)
}

func TestDecodeTypePointer(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	q := r.URL.Query()
	q.Add("val", "v1")
	r.URL.RawQuery = q.Encode()

	type DataType struct {
		Val string `inreq:"query"`
	}

	v, err := DecodeType[*DataType](r)
	require.NoError(t, err)
	require.Equal(t, "v1", v.Val)
}

func TestDecodeTypeDecoder(t *testing.T) {
	type DataType struct {
		Val string `inreq:"query"`
	}

	d := NewTypeDecoder[DataType]()

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	q := r.URL.Query()
	q.Add("val", "v1")
	r.URL.RawQuery = q.Encode()

	v, err := d.Decode(r)
	require.NoError(t, err)
	require.Equal(t, "v1", v.Val)
}

func TestDecodeTypeMapTags(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("val", "x1")
	r.Header.Set("x1", "x2")

	type DataType struct {
		Val string
		X   struct {
			X1 string
		}
	}

	data, err := DecodeType[DataType](r, WithMapTags(map[string]any{
		"Val": "header",
		"X": map[string]any{
			"X1": "header",
		},
	}))
	require.NoError(t, err)
	require.Equal(t, "x1", data.Val)
}

func TestDecodeTypeMapTagsOverrideStructTags(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/?val=x1", nil)
	r.Header.Set("x1", "x2")

	type DataType struct {
		Val string `inreq:"query"`
		X   struct {
			X1 string `inreq:"query"`
		}
	}

	data, err := DecodeType[DataType](r, WithMapTags(map[string]any{
		"X": map[string]any{
			"X1": "header",
		},
	}))
	require.NoError(t, err)
	require.Equal(t, "x1", data.Val)
	require.Equal(t, "x2", data.X.X1)
}
