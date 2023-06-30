package inreq

import "github.com/RangelReale/instruct"

// DecodeContext is the context sent to DecodeOperation.
type DecodeContext interface {
	instruct.DecodeContext
	// PathValue is the function used to extract the path from the request.
	PathValue() PathValue
	// IsBodyDecoded returns whether the body was already decoded.
	IsBodyDecoded() bool
	// DecodedBody signals that the body was decoded.
	DecodedBody()
	// SliceSplitSeparator returns the string used for string-to-array conversions. The default is ",".
	SliceSplitSeparator() string
	// AllowReadBody returns whether the user gave permission to read the request body.
	AllowReadBody() bool
	// EnsureAllQueryUsed returns whether to check if all query parameters were used.
	EnsureAllQueryUsed() bool
	// EnsureAllFormUsed returns whether to check if all form parameters were used.
	EnsureAllFormUsed() bool
}

type decodeContext struct {
	instruct.DefaultDecodeContext
	pathValue           PathValue
	decodedBody         bool
	allowReadBody       bool
	sliceSplitSeparator string
	ensureAllQueryUsed  bool
	ensureAllFormUsed   bool
}

func (d *decodeContext) PathValue() PathValue {
	return d.pathValue
}

func (d *decodeContext) IsBodyDecoded() bool {
	return d.decodedBody
}

func (d *decodeContext) DecodedBody() {
	d.decodedBody = true
}

func (d *decodeContext) AllowReadBody() bool {
	return d.allowReadBody
}

func (d *decodeContext) SliceSplitSeparator() string {
	return d.sliceSplitSeparator
}

func (d *decodeContext) EnsureAllQueryUsed() bool {
	return d.ensureAllQueryUsed
}

func (d *decodeContext) EnsureAllFormUsed() bool {
	return d.ensureAllFormUsed
}
