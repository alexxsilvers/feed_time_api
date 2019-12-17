package predict

import (
	"sync"
	"time"

	"github.com/alexxsilvers/feed_time_api/domain/entity"
)

type predictCache struct {
	ttlInSeconds time.Duration
	data         map[entity.Point]item
	sync.RWMutex
}

type item struct {
	predictTime  int
	timeToRemove int64
}

func newPredictCache(size int, ttlInSeconds int) *predictCache {
	cache := &predictCache{
		ttlInSeconds: time.Duration(ttlInSeconds),
		data:         make(map[entity.Point]item, size),
	}

	ticker := time.NewTicker(time.Second * 20)
	go func() {
		for range ticker.C {
			nowUnixTime := time.Now().Unix()
			cache.Lock()
			for point, item := range cache.data {
				if item.timeToRemove < nowUnixTime {
					delete(cache.data, point)
				}
			}
			cache.Unlock()
		}
	}()

	return cache
}

func (pc *predictCache) set(point *entity.Point, predictTime int) {
	pc.RLock()
	pc.data[*point] = item{
		predictTime:  predictTime,
		timeToRemove: time.Now().Add(time.Second * pc.ttlInSeconds).Unix(),
	}
	pc.RUnlock()
}

func (pc *predictCache) get(point *entity.Point) (int, error) {
	pc.Lock()
	predictTimes, exist := pc.data[*point]
	pc.Unlock()

	if exist {
		return predictTimes.predictTime, nil
	}
	return 0, errNotFoundInCache
}
