package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Browser BrowserConfig `yaml:"browser"`
	Options OptionsConfig `yaml:"options"`
	Server  ServerConfig  `yaml:"server"`
}

type BrowserConfig struct {
	WsURL             string        `yaml:"ws_url" env:"BROWSER_WS_URL" env-required:"true"`
	Referer           string        `yaml:"referer" env-default:"https://google.com"`
	AcceptLanguage    string        `yaml:"accept_language" env-default:"ru-RU,ru;q=0.9"`
	DomStableDuration time.Duration `yaml:"dom_stable_duration" env-default:"2500ms"`
	DomStableDiff     float64       `yaml:"dom_stable_diff" env-default:"0.85"`
	HeadlessMode      bool          `yaml:"headless_mode" env:"BROWSER_HEADLESS_MODE" env-default:"true"`
}

type OptionsConfig struct {
	LoggerTimeFormat string `yaml:"logger_time_format" env:"OPTIONS_LOGGER_TIME_FORMAT" env-default:"02-01-2006 15:04:05"`
}

type ServerConfig struct {
	Env            string        `yaml:"env" env:"SERVER_ENV" env-required:"true"`
	HTTPAddr       string        `yaml:"http_addr" env:"SERVER_HTTP_ADDR" env-required:"true"`
	WbCfg          *WbConfig     `yaml:"wb_config"`
	OzonCfg        *OzonConfig   `yaml:"ozon_config"`
	RequestTimeout time.Duration `yaml:"request_timeout" env:"SERVER_REQUEST_TIMEOUT" env-default:"30s"`
}

type WbConfig struct {
	BaseURL             string `yaml:"base_url" env-required:"true"`
	CloseButtonSelector string `yaml:"close_button_selector" env-required:"true"`
	SearchBarSelector   string `yaml:"search_bar_selector" env-required:"true"`
	ItemsSelector       string `yaml:"items_selector" env-required:"true"`
	LinkSelector        string `yaml:"link_selector" env-required:"true"`
	PriceSelector       string `yaml:"price_selector" env-required:"true"`
	RatingSelector      string `yaml:"rating_selector" env-required:"true"`
	ReviewsSelector     string `yaml:"reviews_selector" env-required:"true"`
}

type OzonConfig struct {
	BaseURL             string `yaml:"base_url" env-required:"true"`
	SearchBarSelector   string `yaml:"search_bar_selector" env-required:"true"`
	ItemsSelector       string `yaml:"items_selector" env-required:"true"`
	LinkSelector        string `yaml:"link_selector" env-required:"true"`
	ProductNameSelector string `yaml:"product_name_selector" env-required:"true"`
	PriceSelector       string `yaml:"price_selector" env-required:"true"`
	RatingSelector      string `yaml:"rating_selector" env-required:"true"`
	ReviewsSelector     string `yaml:"reviews_selector" env-required:"true"`
}

func LoadConfig() (*Config, error) {
	var cfg Config

	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		return nil, fmt.Errorf("CONFIG_PATH not set")
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}



	/*
	if err := cfg.setEnvOptions(); err != nil {
		return nil, fmt.Errorf("set required options: %w", err)
	}
		*/

	return &cfg, nil
}

/*
func (cfg *Config) setEnvOptions() error {
	if env := os.Getenv("SERVER_ENV"); env == "" {
		return fmt.Errorf("SERVER_ENV not set")
	} else {
		cfg.Server.Env = env
	}

	if httpAddr := os.Getenv("SERVER_HTTP_ADDR"); httpAddr == "" {
		return fmt.Errorf("SERVER_HTTP_ADDR not set")
	} else {
		cfg.Server.HTTPAddr = httpAddr
	}

	if wsURL := os.Getenv("BROWSER_WS_URL"); wsURL == "" {
		return fmt.Errorf("BROWSER_WS_URL not set")
	} else {
		cfg.Browser.WsURL = wsURL
	}

	return nil
}
	*/
