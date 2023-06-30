package inreq

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	type DTestEmbed struct {
		H string `inreq:"header"`
		Q string `inreq:"query"`
	}

	type DTest1 struct {
		P string `inreq:"path"`
		Q string `inreq:"query,name=Q1"`
		F string `inreq:"form"`
	}

	type DTestBody struct {
		F1 string
		F2 int
	}

	type DTest struct {
		DTestEmbed
		T1 DTest1    `inreq:"recurse"`
		TB DTestBody `inreq:"body"`
	}

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"F1":"ValueF1","F2":99}`))
	err := r.ParseForm()
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("H", "ValueH")
	q := r.URL.Query()
	q.Add("q", "ValueQ")
	q.Add("Q1", "ValueQ1")
	r.URL.RawQuery = q.Encode()
	r.Form.Add("f", "ValueF")

	data := &DTest{}
	want := &DTest{
		DTestEmbed: DTestEmbed{
			H: "ValueH",
			Q: "ValueQ",
		},
		T1: DTest1{
			P: "ValueP",
			Q: "ValueQ1",
			F: "ValueF",
		},
		TB: DTestBody{
			F1: "ValueF1",
			F2: 99,
		},
	}

	err = Decode(r, data,
		WithPathValue(PathValueFunc(func(r *http.Request, name string) (found bool, value any, err error) {
			if name == "p" {
				return true, "ValueP", err
			}
			return false, nil, nil
		})))
	require.NoError(t, err)
	require.Equal(t, want, data)
}

func TestDecodeEmbed(t *testing.T) {
	type EmbedTestInner struct {
		Val string `inreq:"header"`
	}

	type EmbedTest struct {
		EmbedTestInner
	}

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("val", "x1")

	data := &EmbedTest{}
	want := &EmbedTest{
		EmbedTestInner{Val: "x1"},
	}

	err := Decode(r, data)
	require.NoError(t, err)
	require.Equal(t, want, data)
}

func TestDecodeNonPointer(t *testing.T) {
	type DataType struct {
		Val string `inreq:"header"`
	}

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("val", "x1")

	var data DataType

	err := Decode(r, data)
	var target *InvalidDecodeError
	require.ErrorAs(t, err, &target)
}

func TestDecodeMapTags(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("val", "x1")
	r.Header.Set("x1", "x2")

	type DataType struct {
		Val string
		X   struct {
			X1 string
		}
	}

	var data DataType

	err := Decode(r, &data, WithMapTags(map[string]any{
		"Val": "header",
		"X": map[string]any{
			"X1": "header",
		},
	}))
	require.NoError(t, err)
	require.Equal(t, "x1", data.Val)
}

func TestDecodeMapTagsOverrideStructTags(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/?val=x1", nil)
	r.Header.Set("x1", "x2")

	type DataType struct {
		Val string `inreq:"query"`
		X   struct {
			X1 string `inreq:"query"`
		}
	}

	var data DataType

	err := Decode(r, &data, WithMapTags(map[string]any{
		"X": map[string]any{
			"X1": "header",
		},
	}))
	require.NoError(t, err)
	require.Equal(t, "x1", data.Val)
	require.Equal(t, "x2", data.X.X1)
}

func TestDecoderMapTags(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("val", "x1")
	r.Header.Set("x1", "x2")

	type DataType struct {
		Val string
		X   struct {
			X1 string
		}
	}

	var data DataType

	d := NewDecoder(WithDefaultMapTags(data, map[string]any{
		"Val": "header",
		"X": map[string]any{
			"X1": "header",
		},
	}))

	err := d.Decode(r, &data)
	require.NoError(t, err)
	require.Equal(t, "x1", data.Val)
}

func TestDecoderMapTagsOverrideStructTags(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/?val=x1", nil)
	r.Header.Set("x1", "x2")

	type DataType struct {
		Val string `inreq:"query"`
		X   struct {
			X1 string `inreq:"query"`
		} `inreq:"recurse"`
	}

	var data DataType

	d := NewDecoder()

	err := d.Decode(r, &data, WithMapTags(map[string]any{
		"X": map[string]any{
			"X1": "header",
		},
	}))
	require.NoError(t, err)
	require.Equal(t, "x1", data.Val)
	require.Equal(t, "x2", data.X.X1)
}
