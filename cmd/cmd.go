package cmd

import (
	"flag"
	"fmt"
)

type CMDArgs struct {
	PathConfig string
}

var pathConfig = flag.String("config", "", "config for application")

func ParseArgs() (*CMDArgs, error) {
	flag.Parse()

	if *pathConfig == "" {
		return nil, fmt.Errorf("arguments should not be empty, config: %s", *pathConfig)
	}

	return &CMDArgs{
		PathConfig: *pathConfig,
	}, nil
}
