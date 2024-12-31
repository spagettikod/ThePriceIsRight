package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/spagettikod/thepriceisright"
)

const (
	PriceIsRightEnv = "TPIR"
)

var (
	flagDebug bool = false
)

func parseArgs() (thepriceisright.Config, error) {
	args := flag.Args()

	if flagDebug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	cfg := thepriceisright.NewConfig()

	if len(args) > 2 {
		// fmt.Println("error: too few parameters")
		// fmt.Println()
		// printUsage()
		// os.Exit(2)
		cfg.AreaCode = args[len(args)-2]
		if !slices.Contains(thepriceisright.AreaCodes, cfg.AreaCode) {
			return cfg, fmt.Errorf("area code has invalid value %s, valid values are: %v", cfg.AreaCode, thepriceisright.AreaCodes)
		}
		var err error
		cfg.MaxPrice, err = strconv.ParseFloat(args[len(args)-1], 64)
		if err != nil {
			return cfg, fmt.Errorf("%v is not a valid price", args[len(args)-1])
		}
	}
	return cfg, nil
}

func printUsage() {
	fmt.Println("Usage: tpir [OPTIONS] [area code] [price]")
	fmt.Println("")
	fmt.Println("The Price Is Right calls www.elprisetjustnu.se to check if the price for electricity")
	fmt.Println("is lower or higher than the given price. If lower the command returns 0, if higher")
	fmt.Println("it returns 1.")
	fmt.Println("")
	fmt.Println("   area code     Valid values are SE1, SE2, SE3 or SE4")
	fmt.Println("   price         Price of electricity in SEK per kWh needs to be lower than this to return 0")
	fmt.Println("")
	fmt.Println("Options:")
	flag.PrintDefaults()
}

func init() {
	flag.BoolVar(&flagDebug, "debug", false, "turn on debug output")
}

func main() {
	slog.Debug("Starting up, parsing flags and arguments")
	flag.Parse()

	cfg, err := parseArgs()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(2)
	}

	slog.Debug(fmt.Sprintf("Will evaluate if the electricity price for area code %s is lower than %v SEK/kWh ", cfg.AreaCode, cfg.MaxPrice))
	cache, err := thepriceisright.NewFileCache(cfg.AreaCode)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(2)
	}
	slog.Debug(fmt.Sprintf("Looking up current price using timestamp %v", time.Now().Local()))
	price, err := cache.TodaysPrices().Price(time.Now().Local())
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(2)
	}
	slog.Debug(fmt.Sprintf("Current price for electricity in area code %s is %v SEK/kWh", cfg.AreaCode, price.SekPerKwh))
	if price.SekPerKwh > cfg.MaxPrice {
		slog.Debug(fmt.Sprintf("The Price Is NOT Right! Current electricity price at %v SEK/kWh is higher than the given maximum price at %v SEK/kWh", price.SekPerKwh, cfg.MaxPrice))
		os.Exit(1)
	}
	slog.Debug(fmt.Sprintf("The Price Is Right! Current electricity price at %v SEK/kWh is lower than the given maximum price at %v SEK/kWh", price.SekPerKwh, cfg.MaxPrice))
}
