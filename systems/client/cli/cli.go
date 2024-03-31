package cli

import (
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/aargeee/whwh/systems/client"
)

const NO_URL_PROVIDED = "NO_URL_PROVIDED"

func ParseFlags() string {
	hookUrl := *flag.String("h", NO_URL_PROVIDED, "specify the hook url")
	flag.Parse()
	return hookUrl
}

func ChangeHookForDockerClient(base_cli string) string {
	if os.Getenv("AARGEEE_IS_DOCKER") == "True" {
		return strings.ReplaceAll(base_cli, "localhost", "host.docker.internal")
	}
	return base_cli
}

type CLI struct {
	out io.Writer
}

func NewCLI(out io.Writer) *CLI {
	return &CLI{
		out: out,
	}
}

func (cli *CLI) BeginCLI() {
	// Set up Slog
	if os.Getenv("AARGEEE_TEST") == "True" {
		slog.SetDefault(slog.New(slog.NewTextHandler(cli.out, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})))
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(cli.out, nil)))

	slog.Info("Starting WHWH Client... Thankyou for choosing me!")

	server_url := getServerUrl()
	slog.Debug("Using ", server_url, " for server")

	hook_url := getHookUrl()
	slog.Info("Received Hook Url: " + hook_url)

	c, sid, err := client.ClientConnect(server_url, hook_url)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info(fmt.Sprintf("All data to %s?stream=%s will be forwarded to %s", server_url, sid, hook_url))

	hook_url = SanitizeForDocker(hook_url)

	unsubscribe, err := client.ClientSubscribe(c, sid, hook_url)
	if err != nil {
		log.Fatal(err)
	}
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit
	unsubscribe()
	log.Println("Client Stopped")
}

func getServerUrl() string {
	if os.Getenv("AARGEEE_TEST") != "True" {
		prod_url := os.Getenv("AARGEEE_PROD_URL")
		if prod_url == "" {
			slog.Error("AARGEEE_PROD_URL not provided!")
			os.Exit(70)
			return ""
		}
		return prod_url
	}
	test_url := os.Getenv("AARGEEE_TEST_URL")
	if test_url == "" {
		return "http://localhost:8000"
	}
	return test_url
}

func getHookUrl() string {
	hook_url := ParseFlags()
	if hook_url == NO_URL_PROVIDED {
		hook_url = os.Getenv("AARGEEE_HOOK_URL")
		if hook_url == "" {
			slog.Error("PLEASE PROVIDE A HOOK URL.")
			os.Exit(69)
			return ""
		}
	}
	return hook_url
}

func SanitizeForDocker(url string) string {
	if os.Getenv("AARGEEE_ENV") == "Docker" {
		return strings.ReplaceAll(url, "localhost", "host.docker.internal")
	}
	return url
}
