package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"time"
)

const (
	PriceIsRightEnv = "TPIR"
)

var (
	ErrNotFound      = errors.New("no price found")
	AreaCodes        = []string{"SE1", "SE2", "SE3", "SE4"}
	flagDebug   bool = false
)

func debug(msg string) {
	if flagDebug {
		log.Println(msg)
	}
}

func parseArgs() (string, float64, error) {
	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("error: too few parameters")
		fmt.Println()
		printUsage()
		os.Exit(2)
	}

	code := args[len(args)-2]
	if !slices.Contains(AreaCodes, code) {
		return "", 0, fmt.Errorf("area code has invalid value %s, valid values are: %v", code, AreaCodes)
	}
	maxPrice, err := strconv.ParseFloat(args[len(args)-1], 64)
	if err != nil {
		return "", 0, fmt.Errorf("%v is not a valid price", args[len(args)-1])
	}
	return code, maxPrice, nil
}

func printUsage() {
	fmt.Println("Usage: tpir <area code> <price>")
	fmt.Println(" area code     valid values are SE1, SE2, SE3 or SE4")
	fmt.Println(" price         price of electricity in SEK per kWh needs to be lower than this to return 0")
	fmt.Println("")

	fmt.Println("The Price Is Right calls www.elprisetjustnu.se to check if the price for electricity")
	fmt.Println("is lower or higher than the given price. If lower the command returns 0, if higher")
	fmt.Println("it returns 1.")
	fmt.Println("")
}

func init() {
	flag.BoolVar(&flagDebug, "debug", false, "turn on debug output")
}

func main() {
	debug("Starting up, parsing flags and arguments")
	flag.Parse()

	areaCode, maxPrice, err := parseArgs()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(2)
	}
	debug(fmt.Sprintf("Will evaluate if the electricity price for area code %s is lower than %v SEK/kWh ", areaCode, maxPrice))
	todays, err := Load(areaCode)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(2)
	}
	debug(fmt.Sprintf("Looking up current price using timestamp %v", time.Now().Local()))
	price, err := todays.Price(time.Now().Local())
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(2)
	}
	debug(fmt.Sprintf("Current price for electricity in area code %s is %v SEK/kWh", areaCode, price.SekPerKwh))
	if price.SekPerKwh > maxPrice {
		debug(fmt.Sprintf("The Price Is NOT Right! Current electricity price at %v SEK/kWh is higher than the given maximum price at %v SEK/kWh", price.SekPerKwh, maxPrice))
		os.Exit(1)
	}
	debug(fmt.Sprintf("The Price Is Right! Current electricity price at %v SEK/kWh is lower than the given maximum price at %v SEK/kWh", price.SekPerKwh, maxPrice))
}
