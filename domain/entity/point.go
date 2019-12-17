package entity

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type validationError struct {
	err error
}

func (ve validationError) Error() string {
	return ve.err.Error()
}
func (ve validationError) IsValidationError() bool {
	return true
}

var (
	errInvalidLatitude  = errors.New("latitude must be between -90 and 90")
	errInvalidLongitude = errors.New("longitude must be between -180 and 180")
)

func NewPoint(lat string, lng string) (*Point, error) {
	latFloat, err := strconv.ParseFloat(strings.TrimSpace(lat), 64)
	if err != nil {
		return nil, errors.Wrap(err, "parse latitude")
	}

	lngFloat, err := strconv.ParseFloat(strings.TrimSpace(lng), 64)
	if err != nil {
		return nil, errors.Wrap(validationError{err: err}, "parse longitude")
	}

	return &Point{
		Latitude:  latFloat,
		Longitude: lngFloat,
	}, nil
}

type Point struct {
	Latitude  float64
	Longitude float64
}

func (p *Point) Validate() error {
	if p.Latitude < -90 || p.Latitude > 90 {
		return errInvalidLatitude
	}

	if p.Longitude < -180 || p.Longitude > 180 {
		return errInvalidLongitude
	}

	return nil
}
