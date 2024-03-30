package cli

import (
	"flag"
	"os"
	"strings"
)

func ParseFlags() string {
	hookUrl := *flag.String("h", ChangeHookForDockerClient(os.Getenv("WHWH_HOOK_URL")), "specify the hook url, defaults to localhost:3000")
	flag.Parse()
	return hookUrl
}

func ChangeHookForDockerClient(base_cli string) string {
	if os.Getenv("AARGEEE_IS_DOCKER") == "True" {
		return strings.ReplaceAll(base_cli, "localhost", "host.docker.internal")
	}
	return base_cli
}
