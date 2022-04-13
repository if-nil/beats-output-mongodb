package main

import (
	"github.com/elastic/beats/v7/filebeat/cmd"
	inputs "github.com/elastic/beats/v7/filebeat/input/default-inputs"
	_ "github.com/huawen0327/beats-output-mongodb"
	"os"
)

func main() {
	if err := cmd.Filebeat(inputs.Init, cmd.FilebeatSettings()).Execute(); err != nil {
		os.Exit(1)
	}
}
