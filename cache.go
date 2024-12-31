package thepriceisright

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type Cache interface {
	AreaCode() string
	Expired() bool
	TodaysPrices() TodaysPrices
	Update() error
}

func fetch(areaCode string) (TodaysPrices, error) {
	now := time.Now().Local()
	todays := NewTodaysPrices()

	// cache was not accepted, fetch from REST API
	url := fmt.Sprintf("https://www.elprisetjustnu.se/api/v1/prices/%d/%02d-%02d_%s.json", now.Year(), now.Month(), now.Day(), areaCode)
	slog.Debug(fmt.Sprintf("Fetching new price list from %s", url))
	resp, err := http.Get(url)
	if err != nil {
		return todays, fmt.Errorf("could not fetch daily prices from %s: %w", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		return todays, fmt.Errorf("calling %s responded with status code %v, expected status %v", url, resp.StatusCode, http.StatusOK)
	}
	slog.Debug("Downloaded price list without any errors, trying to read the price list")
	prices := []Price{}
	bites, err := io.ReadAll(resp.Body)
	if err != nil {
		return todays, fmt.Errorf("error while reading response from %s: %w", url, err)
	}
	if err := json.Unmarshal(bites, &prices); err != nil {
		return todays, fmt.Errorf("error while marshaling response from %s: %w", url, err)
	}

	todays.SetPrices(prices)

	slog.Debug("New price list read without errors")

	return todays, nil
}
