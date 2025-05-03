package log

import (
	"context"
	"github.com/lolizeppelin/micro/log/internal"
	"testing"
)

func TestLogging(t *testing.T) {

	ctx := context.Background()

	internal.IgnorePC = false

	Errorf(ctx, "test error 5s %s", "slog")
	Infof(ctx, "test no err 5s %s", "slog")

	//logger.ErrorContext(ctx, "wtf")
}
