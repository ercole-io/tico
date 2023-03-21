package main

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/ercole-io/tico/config"
	"github.com/ercole-io/tico/controller"
	"github.com/fnproject/fdk-go"
)

func main() {
	doc, err := os.ReadFile("config.toml")
	if err != nil {
		panic(err)
	}

	err = toml.Unmarshal(doc, &config.Conf)
	if err != nil {
		panic(err)
	}

	fdk.Handle(fdk.HandlerFunc(controller.ServiceNowHandler))
	fdk.Handle(fdk.HandlerFunc(controller.OracleCloudHandler))
}
