package integration

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/testcontainers/testcontainers-go"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/config"
)

var (
	envPath       = "../../../.env"
	configPathKey = "CONFIG_PATH"
	configPathVal = "../../../configs/config.yaml"

	chromiumPort = "7317/tcp"
	BrowserWSURL = ""

	Cfg *config.Config
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if err := godotenv.Load(envPath); err != nil {
		fmt.Printf("load env: %s", err.Error())
		os.Exit(1)
	}

	if os.Getenv("TEST_INTEGRATION") != "1" {
		fmt.Println("integration tests skipped; set TEST_INTEGRATION=1 to run")
		os.Exit(0)
	}

	if err := os.Setenv(configPathKey, configPathVal); err != nil {
		fmt.Printf("set env: %s", err.Error())
		os.Exit(1)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("load config: %s", err.Error())
		os.Exit(1)
	}

	Cfg = cfg

	req := testcontainers.ContainerRequest{
		Image:        "ghcr.io/go-rod/rod:latest",
		Name:         "testchromium",
		ExposedPorts: []string{chromiumPort},
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Printf("generic container: %s", err)
		os.Exit(1)
	}
	defer terminate(ctx, c)

	host, err := c.Host(ctx)
	if err != nil {
		log.Printf("host: %s", err)
		terminate(ctx, c)
		os.Exit(1)
	}

	port, err := c.MappedPort(ctx, "7317")
	if err != nil {
		log.Printf("mapped port: %s", err)
		terminate(ctx, c)
		os.Exit(1)
	}

	BrowserWSURL = fmt.Sprintf("ws://%s:%s", host, port.Port())
	Cfg.Browser.WsURL = BrowserWSURL

	time.Sleep(time.Second * 3)
	code := m.Run()

	os.Exit(code)
}

func terminate(ctx context.Context, c testcontainers.Container) {
	if c == nil {
		return
	}
	if err := c.Terminate(ctx); err != nil {
		log.Printf("failed terminate container: %v", err)
	}
}
