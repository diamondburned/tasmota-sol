package main

import (
	_ "embed"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/diamondburned/solar"
	"github.com/joho/godotenv"
)

//go:embed .env.example
var defaultEnv string

type Config struct {
	Latitude  float64
	Longitude float64

	Endpoint string

	BrightnessDay   int
	BrightnessNight int

	BulbWarm int
	BulbCold int

	WarmTemperature solar.Temperature
	ColdTemperature solar.Temperature
}

func mustLoadEnv() {
	def, err := godotenv.Unmarshal(defaultEnv)
	if err != nil {
		log.Panicln("BUG: cannot load default .env:", err)
	}

	for k, v := range def {
		os.Setenv(k, v)
	}

	files, _ := filepath.Glob(".env*")

	if err := godotenv.Load(files...); err != nil {
		log.Fatalln("cannot load .env:", err)
	}
}

func mustParseEnvConfig() Config {
	return Config{
		Latitude:        mustParseEnvFloat64("SOL_LATITUDE"),
		Longitude:       mustParseEnvFloat64("SOL_LONGITUDE"),
		Endpoint:        mustEnvString("SOL_ENDPOINT"),
		BrightnessDay:   mustParseEnvInt("SOL_BRIGHTNESS_DAY"),
		BrightnessNight: mustParseEnvInt("SOL_BRIGHTNESS_NIGHT"),
		BulbWarm:        mustParseEnvInt("SOL_BULB_WARM"),
		BulbCold:        mustParseEnvInt("SOL_BULB_COLD"),
		WarmTemperature: solar.Temperature(mustParseEnvFloat64("SOL_WARM_TEMPERATURE")),
		ColdTemperature: solar.Temperature(mustParseEnvFloat64("SOL_COLD_TEMPERATURE")),
	}
}

func mustEnvString(env string) string {
	v := os.Getenv(env)
	if v == "" {
		log.Fatalf("env $%s missing", env)
	}
	return v
}

func mustParseEnvInt(env string) int {
	v, err := strconv.Atoi(os.Getenv(env))
	if err != nil {
		log.Fatalf("env $%s not valid integer: %s", env, err)
	}
	return v
}

func mustParseEnvFloat64(env string) float64 {
	v, err := strconv.ParseFloat(os.Getenv(env), 64)
	if err != nil {
		log.Fatalf("env $%s not valid float: %s", env, err)
	}
	return v
}

func mustParseEnvBool(env string) bool {
	v, err := strconv.ParseBool(os.Getenv(env))
	if err != nil {
		log.Fatalf("env $%s not valid bool: %s", env, err)
	}
	return v
}
