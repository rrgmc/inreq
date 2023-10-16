package resolver

import (
	"github.com/rrgmc/instruct/resolver"
)

// ValueResolver resolves simple types for a Resolver.
// It should NOT handle slices, pointers, or maps.
type ValueResolver resolver.ValueResolver

// TypeValueResolver is a custom type handler for a ValueResolver.
// It should NOT process value using reflection (for performance reasons).
type TypeValueResolver = resolver.TypeValueResolver

// TypeValueResolverReflect is a custom type handler for a ValueResolver.
// It SHOULD process value using reflection.
type TypeValueResolverReflect = resolver.TypeValueResolverReflect
