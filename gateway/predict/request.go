package predict

import "github.com/alexxsilvers/feed_time_api/domain/entity"

type point struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type getRouteTimesRequest struct {
	Target point   `json:"target"`
	Source []point `json:"source"`
}

func newGetRouteTimesRequest(cars []entity.Car, target *entity.Point) getRouteTimesRequest {
	req := getRouteTimesRequest{
		Target: point{
			Latitude:  target.Latitude,
			Longitude: target.Longitude,
		},
		Source: make([]point, 0, len(cars)),
	}

	for _, car := range cars {
		req.Source = append(req.Source, point{
			Latitude:  car.Point.Latitude,
			Longitude: car.Point.Longitude,
		})
	}

	return req
}
