package config

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Mode int

const (
	ModeDev = iota
	ModeProd
)

func (m Mode) String() string {
	switch m {
	case ModeDev:
		return "DEV"
	case ModeProd:
		return "PROD"
	default:
		return ""
	}
}

func (m Mode) MarshalText() ([]byte, error) {
	return []byte(m.String()), nil
}

func (m *Mode) UnmarshalText(text []byte) error {
	switch strings.ToUpper(string(text)) {
	case "DEV":
		*m = ModeDev
	case "PROD":
		*m = ModeProd
	default:
		return errors.New("unknown name")
	}

	return nil
}

type Config struct {
	Mode     Mode
	AppRoot  string
	LogLevel slog.Level

	ServerAddr string
	DSN        string
	HashKey    []byte
	BlockKey   []byte
}

func Load() Config {
	var mode string
	flag.StringVar(&mode, "mode", "DEV", "determine application mode (DEV|PROD)")
	flag.Parse()

	cfg := Config{}
	if err := cfg.Mode.UnmarshalText([]byte(mode)); err != nil {
		fmt.Fprintf(os.Stderr, "invalid argument '%s' for '-mode'\n", mode)
		os.Exit(1)
	}

	switch cfg.Mode {
	case ModeDev:
		cfg.Mode = ModeDev

		if err := godotenv.Load(".env.local"); err != nil && !errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "godotenv: %v\n", err)
			os.Exit(1)
		}
		if err := godotenv.Load(".env"); err != nil {
			fmt.Fprintf(os.Stderr, "godotenv: %v\n", err)
			os.Exit(1)
		}
	case ModeProd:
		cfg.Mode = ModeProd
	}

	if root, ok := os.LookupEnv("APP_ROOT"); ok {
		cfg.AppRoot = root
	} else {
		cfg.AppRoot = "."
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if err := cfg.LogLevel.UnmarshalText([]byte(logLevel)); err != nil {
		fmt.Fprintf(os.Stderr, "invalid value '%s' for 'LOG_LEVEL'\n", logLevel)
		os.Exit(1)
	}

	if addr, ok := os.LookupEnv("SERVER_ADDR"); ok {
		cfg.ServerAddr = addr
	} else {
		cfg.ServerAddr = "127.0.0.1:3000"
	}

	cfg.DSN = os.Getenv("DSN")

	hashKey, err := hex.DecodeString(os.Getenv("HASH_KEY"))
	if err != nil {
		panic(err)
	}
	cfg.HashKey = hashKey

	blockKey, err := hex.DecodeString(os.Getenv("BLOCK_KEY"))
	if err != nil {
		panic(err)
	}
	cfg.BlockKey = blockKey

	return cfg
}
