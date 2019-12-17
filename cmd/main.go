package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pkg/errors"

	"github.com/alexxsilvers/feed_time_api/config"
	"github.com/alexxsilvers/feed_time_api/core"
	"github.com/alexxsilvers/feed_time_api/core/logger"
	"github.com/alexxsilvers/feed_time_api/domain/usecase/get_feed_time"
	"github.com/alexxsilvers/feed_time_api/gateway/cars"
	"github.com/alexxsilvers/feed_time_api/gateway/predict"
)

func main() {
	ctx := context.Background()

	app, err := core.NewApp()
	if err != nil {
		logger.Fatal(ctx, errors.Wrap(err, "create app"))
	}

	carsRepository, err := cars.NewCarsRepository(config.GetCarsApiEndpoint(), config.GetCarsApiTimeoutInSeconds())
	if err != nil {
		logger.Fatal(ctx, errors.WithMessage(err, "init cars repository"))
	}

	predictRepository, err := predict.NewPredictRepository(config.GetPredictApiEndpoint(), config.GetPredictApiTimeoutInSeconds())
	if err != nil {
		logger.Fatal(ctx, errors.WithMessage(err, "init predict repository"))
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.URLFormat)
	r.Get("/feedTime/{lat}/{lng}/", func(writer http.ResponseWriter, request *http.Request) {
		startTS := time.Now()
		defer func() {
			logger.Info(ctx, fmt.Sprintf("%s", time.Since(startTS).String()))
		}()

		usecase := get_feed_time.New(carsRepository, predictRepository)

		minFeedTime, err := usecase.Exec(ctx, chi.URLParam(request, "lat"), chi.URLParam(request, "lng"))
		if err != nil {
			if _, ok := errors.Cause(err).(validationError); ok {
				writer.WriteHeader(http.StatusBadRequest)
				return
			}

			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(writer, minFeedTime)
	})

	app.RegisterPublicServer(r)

	err = app.Run()
	if err != nil {
		logger.Fatal(ctx, errors.Wrap(err, "run app"))
	}
}

type validationError interface {
	IsValidationError() bool
}
