package thepriceisright

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"sync"
	"time"
)

var (
	ErrCacheFileNotFound = errors.New("price list cache file not found")
)

type FileCache struct {
	mu          *sync.Mutex
	pricesCache TodaysPrices
	areaCode    string
}

func NewFileCache(areaCode string) (fc *FileCache, err error) {
	fc = &FileCache{
		mu:       &sync.Mutex{},
		areaCode: areaCode,
	}
	fc.pricesCache, err = fc.loadCache()
	if fc.Expired() {
		err = fc.Update()
	}
	return
}

func (fc FileCache) AreaCode() string {
	return fc.areaCode
}

func (fc FileCache) Expired() bool {
	if !fc.pricesCache.IsValid() {
		return true
	}
	return fc.pricesCache.IsExpired(time.Now())
}

func (mc FileCache) TodaysPrices() TodaysPrices {
	return mc.pricesCache
}

func (fc FileCache) Update() error {
	tp, err := fetch(fc.areaCode)
	if err != nil {
		return err
	}
	fc.mu.Lock()
	defer fc.mu.Unlock()
	bites, err := json.Marshal(tp)
	if err != nil {
		return fmt.Errorf("error while marshaling prices")
	}
	cachePath, err := fc.cachePath()
	if err != nil {
		return err
	}
	slog.Debug(fmt.Sprintf("Saving new price list cache file to %s", cachePath))
	return os.WriteFile(cachePath, bites, 0660)
}

func (fc FileCache) cachePath() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	dir = path.Join(dir, "thepriceisright")
	if err := os.MkdirAll(dir, 0750); err != nil {
		return "", err
	}
	return path.Join(dir, fc.areaCode+"_cache.json"), nil
}

func (fc FileCache) loadCache() (TodaysPrices, error) {
	todays := TodaysPrices{Prices: []Price{}}
	cachePath, err := fc.cachePath()
	if err != nil {
		return todays, err
	}
	slog.Debug(fmt.Sprintf("Looking for price list cache file at %s", cachePath))
	b, err := os.ReadFile(cachePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return todays, fmt.Errorf("error while loading temporary cache file: %w", err)
		}
		slog.Debug(fmt.Sprintf("Price list cache file not found at %s", cachePath))
		return todays, ErrCacheFileNotFound
	}
	slog.Debug("Reading price list")
	if err := json.Unmarshal(b, &todays.Prices); err != nil {
		return todays, fmt.Errorf("error while marshaling temporary cache file: %w", err)
	}
	return todays, nil
}
