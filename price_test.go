package thepriceisright

import (
	"testing"
	"time"
)

const (
	prices = `[{"SEK_per_kWh":0.99498,"EUR_per_kWh":0.08883,"EXR":11.200938,"time_start":"2024-01-05T00:00:00+01:00","time_end":"2024-01-05T01:00:00+01:00"},{"SEK_per_kWh":0.93069,"EUR_per_kWh":0.08309,"EXR":11.200938,"time_start":"2024-01-05T01:00:00+01:00","time_end":"2024-01-05T02:00:00+01:00"},{"SEK_per_kWh":0.89608,"EUR_per_kWh":0.08,"EXR":11.200938,"time_start":"2024-01-05T02:00:00+01:00","time_end":"2024-01-05T03:00:00+01:00"},{"SEK_per_kWh":0.89798,"EUR_per_kWh":0.08017,"EXR":11.200938,"time_start":"2024-01-05T03:00:00+01:00","time_end":"2024-01-05T04:00:00+01:00"},{"SEK_per_kWh":0.92397,"EUR_per_kWh":0.08249,"EXR":11.200938,"time_start":"2024-01-05T04:00:00+01:00","time_end":"2024-01-05T05:00:00+01:00"},{"SEK_per_kWh":0.98725,"EUR_per_kWh":0.08814,"EXR":11.200938,"time_start":"2024-01-05T05:00:00+01:00","time_end":"2024-01-05T06:00:00+01:00"},{"SEK_per_kWh":1.22975,"EUR_per_kWh":0.10979,"EXR":11.200938,"time_start":"2024-01-05T06:00:00+01:00","time_end":"2024-01-05T07:00:00+01:00"},{"SEK_per_kWh":1.84827,"EUR_per_kWh":0.16501,"EXR":11.200938,"time_start":"2024-01-05T07:00:00+01:00","time_end":"2024-01-05T08:00:00+01:00"},{"SEK_per_kWh":3.34919,"EUR_per_kWh":0.29901,"EXR":11.200938,"time_start":"2024-01-05T08:00:00+01:00","time_end":"2024-01-05T09:00:00+01:00"},{"SEK_per_kWh":2.79038,"EUR_per_kWh":0.24912,"EXR":11.200938,"time_start":"2024-01-05T09:00:00+01:00","time_end":"2024-01-05T10:00:00+01:00"},{"SEK_per_kWh":2.86587,"EUR_per_kWh":0.25586,"EXR":11.200938,"time_start":"2024-01-05T10:00:00+01:00","time_end":"2024-01-05T11:00:00+01:00"},{"SEK_per_kWh":2.36463,"EUR_per_kWh":0.21111,"EXR":11.200938,"time_start":"2024-01-05T11:00:00+01:00","time_end":"2024-01-05T12:00:00+01:00"},{"SEK_per_kWh":2.27693,"EUR_per_kWh":0.20328,"EXR":11.200938,"time_start":"2024-01-05T12:00:00+01:00","time_end":"2024-01-05T13:00:00+01:00"},{"SEK_per_kWh":2.20636,"EUR_per_kWh":0.19698,"EXR":11.200938,"time_start":"2024-01-05T13:00:00+01:00","time_end":"2024-01-05T14:00:00+01:00"},{"SEK_per_kWh":2.21622,"EUR_per_kWh":0.19786,"EXR":11.200938,"time_start":"2024-01-05T14:00:00+01:00","time_end":"2024-01-05T15:00:00+01:00"},{"SEK_per_kWh":3.35961,"EUR_per_kWh":0.29994,"EXR":11.200938,"time_start":"2024-01-05T15:00:00+01:00","time_end":"2024-01-05T16:00:00+01:00"},{"SEK_per_kWh":5.03639,"EUR_per_kWh":0.44964,"EXR":11.200938,"time_start":"2024-01-05T16:00:00+01:00","time_end":"2024-01-05T17:00:00+01:00"},{"SEK_per_kWh":5.89449,"EUR_per_kWh":0.52625,"EXR":11.200938,"time_start":"2024-01-05T17:00:00+01:00","time_end":"2024-01-05T18:00:00+01:00"},{"SEK_per_kWh":3.91977,"EUR_per_kWh":0.34995,"EXR":11.200938,"time_start":"2024-01-05T18:00:00+01:00","time_end":"2024-01-05T19:00:00+01:00"},{"SEK_per_kWh":2.03756,"EUR_per_kWh":0.18191,"EXR":11.200938,"time_start":"2024-01-05T19:00:00+01:00","time_end":"2024-01-05T20:00:00+01:00"},{"SEK_per_kWh":1.67566,"EUR_per_kWh":0.1496,"EXR":11.200938,"time_start":"2024-01-05T20:00:00+01:00","time_end":"2024-01-05T21:00:00+01:00"},{"SEK_per_kWh":1.26369,"EUR_per_kWh":0.11282,"EXR":11.200938,"time_start":"2024-01-05T21:00:00+01:00","time_end":"2024-01-05T22:00:00+01:00"},{"SEK_per_kWh":1.12121,"EUR_per_kWh":0.1001,"EXR":11.200938,"time_start":"2024-01-05T22:00:00+01:00","time_end":"2024-01-05T23:00:00+01:00"},{"SEK_per_kWh":1.04023,"EUR_per_kWh":0.09287,"EXR":11.200938,"time_start":"2024-01-05T23:00:00+01:00","time_end":"2024-01-06T00:00:00+01:00"}]`
)

func TestPrice(t *testing.T) {
	type TestCase struct {
		date          time.Time
		expectedPrice float64
		expectedError error
	}
	testCases := []TestCase{
		{
			date:          time.Date(2024, time.January, 5, 10, 59, 59, 0, time.Local),
			expectedPrice: 2.86587,
			expectedError: nil,
		},
		{
			date:          time.Date(2024, time.January, 5, 11, 0, 0, 0, time.Local),
			expectedPrice: 2.86587,
			expectedError: nil,
		},
		{
			date:          time.Date(2000, time.January, 5, 11, 0, 0, 0, time.Local),
			expectedPrice: 0,
			expectedError: ErrNotFound,
		},
	}

	todays, err := parse([]byte(prices))
	if err != nil {
		t.Fatalf("unexpected error occurred parsing json: %s", err)
	}
	for _, tc := range testCases {
		price, err := todays.Price(tc.date)
		if err != tc.expectedError {
			t.Fatalf("%v, unexpected error occurred: %s", tc.date, err)
		}
		if price.SekPerKwh != tc.expectedPrice {
			t.Fatalf("invalid price, expected %v but got %v", tc.expectedPrice, price)
		}
	}
}

func TestIsExpired(t *testing.T) {
	todays, err := parse([]byte(prices))
	if err != nil {
		t.Fatalf("unexpected error occurred parsing json: %s", err)
	}
	notExpired := time.Date(2024, time.January, 5, 23, 59, 59, 0, time.Local)
	if todays.IsExpired(notExpired) {
		t.Fatalf("expected not to be expired but was")
	}
	expired := time.Date(2024, time.January, 6, 0, 0, 0, 0, time.Local)
	if !todays.IsExpired(expired) {
		t.Fatalf("expected to be expired but was not")
	}
	expired = time.Date(2024, time.January, 7, 0, 0, 0, 0, time.Local)
	if !todays.IsExpired(expired) {
		t.Fatalf("expected to be expired but was not")
	}
}
