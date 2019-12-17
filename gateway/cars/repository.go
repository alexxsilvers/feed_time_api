package cars

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"

	"github.com/alexxsilvers/feed_time_api/domain/entity"
)

var (
	errEmptyEndpoint   = errors.New("empty endpoint")
	errInvalidEndpoint = errors.New("invalid endpoint")
	errTimeoutIsZero   = errors.New("timeout is zero")
	errBadRequest      = errors.New("bad request")
	errServerError     = errors.New("server error")
	errInvalidLimit    = errors.New("limit must be from 1 to 100")
)

type repo struct {
	httpClient *http.Client
	endpoint   string
}

func NewCarsRepository(endpoint string, timeoutInSeconds int) (*repo, error) {
	if endpoint == "" {
		return nil, errEmptyEndpoint
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "parse endpoint url")
	}

	if u.Host == "" || u.Scheme == "" {
		return nil, errInvalidEndpoint
	}

	endpoint = u.Scheme + "://" + u.Host + u.Path

	if timeoutInSeconds == 0 {
		return nil, errTimeoutIsZero
	}

	r := &repo{
		endpoint: endpoint,
		httpClient: &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				DisableKeepAlives:     false,
				DisableCompression:    false,
			},
			Timeout: time.Duration(timeoutInSeconds) * time.Second,
		},
	}

	return r, nil
}

func (r *repo) GetCars(ctx context.Context, point *entity.Point, limit int) ([]entity.Car, error) {
	if limit < 1 || limit > 100 {
		return nil, errInvalidLimit
	}

	requestURL := fmt.Sprintf("%s?lat=%.8f&lng=%.8f&limit=%d", r.endpoint, point.Latitude, point.Longitude, limit)
	httpReq, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "create http request")
	}

	resp, err := r.httpClient.Do(httpReq)
	// When you get a redirection error, response will not be nil but there’ll be an error
	if resp != nil {
		defer func() {
			errBodyClose := resp.Body.Close()
			if errBodyClose != nil {
				err = errors.Wrap(errBodyClose, "close response body")
			}
		}()
	}
	if err != nil {
		return nil, errors.Wrap(err, "do request")
	}

	if resp == nil {
		return nil, errors.New("empty response")
	}

	if resp.StatusCode == http.StatusBadRequest {
		return nil, errBadRequest
	} else if resp.StatusCode == http.StatusInternalServerError {
		return nil, errServerError
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("server error. code: %d message: %s", resp.StatusCode, resp.Status)
	}

	rawResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read resp body")
	}

	var cars []Car
	err = json.Unmarshal(rawResp, &cars)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal response")
	}

	domainCars := make([]entity.Car, 0, len(cars))
	for _, car := range cars {
		domainCars = append(domainCars, mapCarToDomainCar(&car))
	}

	return domainCars, nil
}
