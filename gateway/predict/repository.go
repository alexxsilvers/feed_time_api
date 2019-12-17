package predict

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/pkg/errors"

	"github.com/alexxsilvers/feed_time_api/domain/entity"
	"github.com/alexxsilvers/feed_time_api/gateway"
)

var (
	errEmptyEndpoint   = errors.New("empty endpoint")
	errInvalidEndpoint = errors.New("invalid endpoint")
	errTimeoutIsZero   = errors.New("timeout is zero")
	errBadRequest      = errors.New("bad request")
	errServerError     = errors.New("server error")
	errNotFoundInCache = errors.New("not found in cache")
)

type repo struct {
	endpoint   string
	cache      *predictCache
	httpClient *http.Client
}

func NewPredictRepository(endpoint string, timeoutInSeconds int) (*repo, error) {
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
		cache:    newPredictCache(1000, 20),
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

func (r *repo) GetRouteTimeFromCache(target *entity.Point) (int, error) {
	return r.cache.get(target)
}

func (r *repo) GetRouteTime(ctx context.Context, cars []entity.Car, target *entity.Point) (int, error) {
	req := newGetRouteTimesRequest(cars, target)

	buf := gateway.NewBuffer()
	defer buf.Free()

	err := json.NewEncoder(buf).Encode(&req)
	if err != nil {
		return 0, errors.Wrap(err, "marshal request")
	}

	httpReq, err := http.NewRequest(http.MethodPost, r.endpoint, buf)
	if err != nil {
		return 0, errors.Wrap(err, "create http request")
	}
	httpReq.Header.Add("Content-Type", "application/json")

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
		return 0, errors.Wrap(err, "do request")
	}

	if resp == nil {
		return 0, errors.New("empty response")
	}

	if resp.StatusCode == http.StatusBadRequest {
		return 0, errBadRequest
	} else if resp.StatusCode == http.StatusInternalServerError {
		return 0, errServerError
	} else if resp.StatusCode != http.StatusOK {
		return 0, errors.Errorf("server error. code: %d message: %s", resp.StatusCode, resp.Status)
	}

	rawResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, errors.Wrap(err, "read resp body")
	}

	var times []int
	err = json.Unmarshal(rawResp, &times)
	if err != nil {
		return 0, errors.Wrap(err, "unmarshal response")
	}

	sort.Ints(times)

	r.cache.set(target, times[0])

	return times[0], nil
}
