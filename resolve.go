package inreq

import "github.com/RangelReale/instruct"

// Resolver converts strings to the type of the struct field.
type Resolver = instruct.Resolver

// DefaultResolve resolve the string value to the proper type and return the value.
var DefaultResolve = instruct.DefaultResolve
