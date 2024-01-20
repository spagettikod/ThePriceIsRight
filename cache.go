package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
)

var (
	ErrCacheFileNotFound = errors.New("price list cache file not found")
)

func cachePath(areaCode string) (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	dir = path.Join(dir, "thepriceisright")
	if err := os.MkdirAll(dir, 0750); err != nil {
		return "", err
	}
	return path.Join(dir, areaCode+"_cache.json"), nil
}

func loadCache(areaCode string) (TodaysPrices, error) {
	todays := TodaysPrices{Prices: []Price{}}
	cachePath, err := cachePath(areaCode)
	if err != nil {
		return todays, err
	}
	debug(fmt.Sprintf("Looking for price list cache file at %s", cachePath))
	b, err := os.ReadFile(cachePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return todays, fmt.Errorf("error while loading temporary cache file: %w", err)
		}
		debug(fmt.Sprintf("Price list cache file not found at %s", cachePath))
		return todays, ErrCacheFileNotFound
	}
	debug("Reading price list")
	if err := json.Unmarshal(b, &todays.Prices); err != nil {
		return todays, fmt.Errorf("error while marshaling temporary cache file: %w", err)
	}
	todays.IsCached = true
	return todays, nil
}
