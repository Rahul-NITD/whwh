package cli

import "flag"

func ParseFlags() string {
	hookUrl := *flag.String("h", "http://localhost:3000", "specify the hook url, defaults to localhost:3000")
	flag.Parse()
	return hookUrl
}
