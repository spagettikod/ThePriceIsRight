package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"slices"
	"strconv"
	"time"
)

const (
	PriceIsRightEnv = "TPIR"
)

var (
	ErrNotFound = errors.New("no price found")
	AreaCodes   = []string{"SE1", "SE2", "SE3", "SE4"}
)

type Price struct {
	SekPerKwh    float64   `json:"SEK_per_kWh"`
	EurPerKwh    float64   `json:"EUR_per_kWh"`
	ExchangeRate float64   `json:"EXR"`
	Start        time.Time `json:"time_start"`
	End          time.Time `json:"time_end"`
}

type TodaysPrices struct {
	Prices []Price
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

func cachePath() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	dir = path.Join(dir, "thepriceisright")
	if err := os.MkdirAll(dir, 0750); err != nil {
		return "", err
	}
	return path.Join(dir, "cache.json"), nil
}

func loadCache() (TodaysPrices, error) {
	todays := TodaysPrices{Prices: []Price{}}
	cachePath, err := cachePath()
	if err != nil {
		return todays, err
	}
	b, err := os.ReadFile(cachePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return todays, fmt.Errorf("error while loading temporary cache file: %w", err)
		}
		return todays, nil
	}
	if err := json.Unmarshal(b, &todays.Prices); err != nil {
		return todays, fmt.Errorf("error while marshaling temporary cache file: %w", err)
	}
	return todays, nil
}

func load(class string) (TodaysPrices, error) {
	now := time.Now().Local()
	todays, err := loadCache()
	if err != nil {
		return todays, fmt.Errorf("error while trying to load from cache: %w", err)
	}

	// if prices in cache is valid and has not expired we return the cache
	if todays.IsValid() && !todays.IsExpired(now) {
		return todays, nil
	}

	// cache was not accepted, fetch from REST API
	url := fmt.Sprintf("https://www.elprisetjustnu.se/api/v1/prices/%d/%02d-%02d_%s.json", now.Year(), now.Month(), now.Day(), class)
	resp, err := http.Get(url)
	if err != nil {
		return todays, fmt.Errorf("could not fetch daily prices from %s: %w", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		return todays, fmt.Errorf("calling %s responded with status code %v, expected status %v", url, resp.StatusCode, http.StatusOK)
	}
	bites, err := io.ReadAll(resp.Body)
	if err != nil {
		return todays, fmt.Errorf("error while reading response from %s: %w", url, err)
	}
	if err := json.Unmarshal(bites, &todays.Prices); err != nil {
		return todays, fmt.Errorf("error while marshaling response from %s: %w", url, err)
	}

	// save fetched prices as cache
	cachePath, err := cachePath()
	if err != nil {
		return todays, err
	}
	os.WriteFile(cachePath, bites, 0660)

	return todays, nil
}

func parseArgs() (string, float64, error) {
	args := os.Args
	if len(args) < 3 {
		fmt.Println("error: too few parameters")
		fmt.Println()
		printUsage()
		os.Exit(2)
	}

	code := args[len(args)-2]
	if !slices.Contains(AreaCodes, code) {
		return "", 0, fmt.Errorf("area code has invalid value %s, valid values are: %v", code, AreaCodes)
	}
	maxPrice, _ := strconv.ParseFloat(args[len(args)-1], 64)
	if maxPrice <= 0 {
		return "", 0, errors.New("price must have a valid value")
	}
	return code, maxPrice, nil
}

func printUsage() {
	fmt.Println("The Price Is Right calls www.elprisetjustnu.se to check if the price for electricity")
	fmt.Println("is lower or higher than the given price. If lower the command returns 0, if higher")
	fmt.Println("it returns 1.")
	fmt.Println("")
	fmt.Println("Usage: tpir <area code> <price>")
	fmt.Println(" area code     valid values are SE1, SE2, SE3 or SE4")
	fmt.Println(" price         price of electricity in SEK per kWh needs to be lower than this to return 0")
	fmt.Println("")
}

func main() {
	areaCode, maxPrice, err := parseArgs()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(3)
	}
	todays, err := load(areaCode)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(3)
	}
	price, err := todays.Price(time.Now().Local())
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(3)
	}

	if price.SekPerKwh > maxPrice {
		os.Exit(1)
	}
}
