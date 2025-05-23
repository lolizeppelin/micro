package transport

import (
	"context"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

type metadataKey struct{}

var (
	eng = cases.Title(language.English)
)

// Metadata is our way of representing request headers internally.
// They're used at the RPC level and translate back and forth
// from Transport headers.
type Metadata map[string]string

func (md Metadata) Get(key string) (string, bool) {
	// attempt to get as is
	val, ok := md[key]
	if ok {
		return val, ok
	}

	// attempt to get lower case
	val, ok = md[strings.Title(key)]
	return val, ok
}

func (md Metadata) Set(key, val string) {
	md[key] = val
}

func (md Metadata) Delete(key string) {
	// delete key as-is
	delete(md, key)
	// delete also Title key
	delete(md, strings.Title(key))
}

// MetadataCopy makes a copy of the metadata.
func MetadataCopy(md Metadata) Metadata {
	cmd := make(Metadata, len(md))
	for k, v := range md {
		cmd[k] = v
	}
	return cmd
}

// ContextDelete key from metadata.
func ContextDelete(ctx context.Context, k string) context.Context {
	return ContextSet(ctx, k, "")
}

// ContextSet add key with val to metadata.
func ContextSet(ctx context.Context, k, v string) context.Context {
	md, ok := ctx.Value(metadataKey{}).(Metadata)
	if !ok {
		md = make(Metadata)
	}
	if v == "" {
		delete(md, k)
	} else {
		md[k] = v
	}
	if ok {
		return ctx
	}
	return context.WithValue(ctx, metadataKey{}, md)
}

// ContextGet returns a single value from metadata in the context.
func ContextGet(ctx context.Context, key string) (string, bool) {
	md, ok := FromContext(ctx)
	if !ok {
		return "", ok
	}
	// attempt to get as is
	val, ok := md[key]
	if ok {
		return val, ok
	}

	// attempt to get lower case
	val, ok = md[strings.Title(key)]

	return val, ok
}

// FromContext returns copied metadata from the given context.
func FromContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(metadataKey{}).(Metadata)
	if !ok {
		return nil, ok
	}

	// capitalise all values
	newMD := make(Metadata, len(md))
	for k, v := range md {
		newMD[eng.String(k)] = v
	}

	return newMD, ok
}

func CopyFromContext(ctx context.Context) map[string]string {
	md, ok := FromContext(ctx)
	if !ok {
		return make(map[string]string)
	}
	return md
}

// NewContext creates a new context with the given metadata.
func NewContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, metadataKey{}, md)
}

// MergeContext merges metadata to existing metadata, overwriting if specified.
func MergeContext(ctx context.Context, patchMd Metadata, overwrite bool) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	md, _ := ctx.Value(metadataKey{}).(Metadata)
	cmd := make(Metadata, len(md))
	for k, v := range md {
		cmd[k] = v
	}
	for k, v := range patchMd {
		if _, ok := cmd[k]; ok && !overwrite {
			// skip
		} else if v != "" {
			cmd[k] = v
		} else {
			delete(cmd, k)
		}
	}
	return context.WithValue(ctx, metadataKey{}, cmd)
}
