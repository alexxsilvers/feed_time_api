package cars

import "github.com/alexxsilvers/feed_time_api/domain/entity"

type Car struct {
	ID        int64   `json:"id"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

func mapCarToDomainCar(car *Car) entity.Car {
	return entity.Car{
		ID: car.ID,
		Point: entity.Point{
			Latitude:  car.Latitude,
			Longitude: car.Longitude,
		},
	}
}
