package transport

import (
	"context"
	"fmt"
	"testing"
)

func TestMD(t *testing.T) {

	ctx := context.Background()

	ctx = ContextSet(ctx, "a", "b")
	ctx = ContextSet(ctx, "c", "d")

	md, _ := FromContext(ctx)
	fmt.Printf("%v", md)
}
