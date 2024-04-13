package client

import (
	"context"
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/utils"
	"time"
)

type BackoffFunc func(ctx context.Context, req micro.Request, attempts int) (time.Duration, error)

func exponentialBackoff(ctx context.Context, req micro.Request, attempts int) (time.Duration, error) {
	if attempts > 13 {
		return 2 * time.Minute, nil
	}
	return utils.BackoffDelay(attempts), nil
}
