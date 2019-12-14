package main

import (
	"context"

	"github.com/pkg/errors"

	"github.com/alexxsilvers/feed_time_api/core"
	"github.com/alexxsilvers/feed_time_api/core/logger"
)

func main() {
	ctx := context.Background()

	app, err := core.NewApp()
	if err != nil {
		logger.Fatal(ctx, errors.Wrap(err, "create app"))
	}

	err = app.Run()
	if err != nil {
		logger.Fatal(ctx, errors.Wrap(err, "run app"))
	}
}
