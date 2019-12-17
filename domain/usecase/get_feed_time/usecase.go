package get_feed_time

import (
	"context"

	"github.com/pkg/errors"

	"github.com/alexxsilvers/feed_time_api/domain/entity"
)

type CarsRepository interface {
	GetCars(ctx context.Context, point *entity.Point, limit int) ([]entity.Car, error)
}

type PredictRepository interface {
	GetRouteTimeFromCache(target *entity.Point) (int, error)
	GetRouteTime(ctx context.Context, cars []entity.Car, target *entity.Point) (int, error)
}

type usecase struct {
	carsRepo    CarsRepository
	predictRepo PredictRepository
}

func New(carsRepo CarsRepository, predictRepo PredictRepository) *usecase {
	return &usecase{
		carsRepo:    carsRepo,
		predictRepo: predictRepo,
	}
}

func (uc *usecase) Exec(ctx context.Context, lat string, lng string) (int, error) {
	point, err := entity.NewPoint(lat, lng)
	if err != nil {
		return 0, errors.Wrap(err, "parse point from string params")
	}

	err = point.Validate()
	if err != nil {
		return 0, errors.Wrap(err, "validate incoming point")
	}

	// try to find in predict cache
	time, err := uc.predictRepo.GetRouteTimeFromCache(point)
	if err == nil {
		return time, nil
	}

	cars, err := uc.carsRepo.GetCars(ctx, point, 100)
	if err != nil {
		return 0, errors.Wrap(err, "get cars around a given point")
	}

	time, err = uc.predictRepo.GetRouteTime(ctx, cars, point)
	if err != nil {
		return 0, errors.Wrap(err, "get predicts around a given point")
	}

	return time, nil
}
