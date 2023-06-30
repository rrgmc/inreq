package inreq

import "github.com/RangelReale/instruct"

// map_tags.go

// MapTags is an alternative to struct tags, and can be used to override them.
type MapTags = instruct.MapTags

// error.go

type ValuesNotUsedError = instruct.ValuesNotUsedError

type InvalidDecodeError = instruct.InvalidDecodeError
