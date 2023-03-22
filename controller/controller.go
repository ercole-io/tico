package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ercole-io/tico/api"
	"github.com/ercole-io/tico/config"
	"github.com/ercole-io/tico/service"
)

func Handler(ctx context.Context, in io.Reader, out io.Writer) {
	c := api.New(config.Conf.ServiceNow.URL, config.Conf.ServiceNow.Username, config.Conf.ServiceNow.Password)

	serviceNowResult, err := c.GetServiceNowResult("cmdb_ci_business_app")
	if err != nil {
		println(err)
	}

	oracleCloud := service.ExampleSearchResources()

	json.NewEncoder(out).Encode(fmt.Sprintf("serviceNow objs:%v\nOracle Cloud objs:%v", len(serviceNowResult.Result), len(oracleCloud)))

}
