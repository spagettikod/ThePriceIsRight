package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Price struct {
	SekPerKwh    float64   `json:"SEK_per_kWh"`
	EurPerKwh    float64   `json:"EUR_per_kWh"`
	ExchangeRate float64   `json:"EXR"`
	Start        time.Time `json:"time_start"`
	End          time.Time `json:"time_end"`
}

type TodaysPrices struct {
	Prices   []Price
	IsCached bool `json:"-"`
}

func (tp TodaysPrices) Price(timestamp time.Time) (Price, error) {
	for _, tp := range tp.Prices {
		if (timestamp.After(tp.Start) || timestamp == tp.Start) && (timestamp.Before(tp.End) || timestamp == tp.End) {
			return tp, nil
		}
	}
	return Price{}, ErrNotFound
}

func (tp TodaysPrices) IsExpired(timestamp time.Time) bool {
	if tp.IsValid() {
		lastPrice := tp.Prices[len(tp.Prices)-1]
		return timestamp.After(lastPrice.End) || timestamp == lastPrice.End
	}
	return true
}

func (tp TodaysPrices) IsValid() bool {
	return len(tp.Prices) == 24
}

func Load(areaCode string) (TodaysPrices, error) {
	now := time.Now().Local()
	todays, err := loadCache(areaCode)
	if err != nil && err != ErrCacheFileNotFound {
		return todays, fmt.Errorf("error while trying to load from cache: %w", err)
	}

	// if prices in cache is valid and has not expired we return the cache
	if todays.IsValid() && !todays.IsExpired(now) {
		debug("Found valid, and current, price list cache file")
		return todays, nil
	}
	if todays.IsCached {
		if todays.IsExpired(now) {
			debug("Cache found but has expired, will download a new one")
		}
		if !todays.IsValid() {
			debug("Cache found but it was invalid, will download a new one")
		}
	}

	// cache was not accepted, fetch from REST API
	url := fmt.Sprintf("https://www.elprisetjustnu.se/api/v1/prices/%d/%02d-%02d_%s.json", now.Year(), now.Month(), now.Day(), areaCode)
	debug(fmt.Sprintf("Fetching new price list from %s", url))
	resp, err := http.Get(url)
	if err != nil {
		return todays, fmt.Errorf("could not fetch daily prices from %s: %w", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		return todays, fmt.Errorf("calling %s responded with status code %v, expected status %v", url, resp.StatusCode, http.StatusOK)
	}
	debug("Downloaded price list without any errors, trying to read the price list")
	bites, err := io.ReadAll(resp.Body)
	if err != nil {
		return todays, fmt.Errorf("error while reading response from %s: %w", url, err)
	}
	if err := json.Unmarshal(bites, &todays.Prices); err != nil {
		return todays, fmt.Errorf("error while marshaling response from %s: %w", url, err)
	}

	debug("New price list read without errors")
	// save fetched prices as cache
	cachePath, err := cachePath(areaCode)
	if err != nil {
		return todays, err
	}
	debug(fmt.Sprintf("Saving new price list cache file to %s", cachePath))
	os.WriteFile(cachePath, bites, 0660)

	return todays, nil
}
