package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/diamondburned/solar"
)

func main() {
	mustLoadEnv()
	cfg := mustParseEnvConfig()

	tick := time.Tick(time.Minute)
	now := time.Now()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var lastTemp solar.Temperature

	// TODO: optimize.
	for {
		temperature, sun := solar.CalculateTemperature(
			now,
			cfg.Latitude, cfg.Longitude,
			cfg.WarmTemperature, cfg.ColdTemperature,
		)

		if lastTemp == temperature {
			// Same results. Skip.
			continue
		}

		// Store the results.
		lastTemp = temperature

		// Map to color. Cold temperature is assumed to be higher than warm:
		// 6500K > 4000K, and 500 > 150.
		color := scaleInt(
			float64(temperature),
			float64(cfg.WarmTemperature), float64(cfg.ColdTemperature),
			cfg.BulbWarm, cfg.BulbCold,
		)

		v := url.Values{
			"t0": {strconv.Itoa(color)},
		}

		log.Println("changing bulb color to", color)

		if sun.IsSetting(now) {
			// Calculate the brightness based on how far we've transitioned
			// using the color temperature.
			brightness := scaleInt(
				float64(temperature),
				float64(cfg.WarmTemperature), float64(cfg.ColdTemperature),
				0, 100,
			)
			v["m"] = []string{"1"}
			v["d0"] = []string{strconv.Itoa(brightness)}
			log.Println("changing bulb brightness to", brightness)
		}

		if err := send(cfg.Endpoint, v); err != nil {
			log.Println("cannot update bulb:", err)
		}

		select {
		case <-ctx.Done():
			return
		case now = <-tick:
			continue
		}
	}
}

// send sends a HTTP GET request.
func send(url string, v url.Values) error {
	r, err := http.Get(url + "?" + v.Encode())
	if err != nil {
		return err
	}
	r.Body.Close()
	if r.StatusCode < 200 || r.StatusCode > 299 {
		return fmt.Errorf("unexpected status %s", r.Status)
	}
	return nil
}

// scaleInt is the rounded int version of scale.
func scaleInt(v, minv, maxv float64, minr, maxr int) int {
	return int(math.Round(scale(v, minv, maxv, float64(minr), float64(maxr))))
}

// scale scales v that is within [minv, maxv] to be within [minr, maxr].
func scale(v, minv, maxv, minr, maxr float64) float64 {
	return minr + ((v - minv) / (maxv - minv) * (maxr - minr))
}
