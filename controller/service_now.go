package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ercole-io/tico/api"
	"github.com/ercole-io/tico/config"
)

func ServiceNowHandler(ctx context.Context, in io.Reader, out io.Writer) {
	c := api.New(config.Conf.ServiceNow.URL, config.Conf.ServiceNow.Username, config.Conf.ServiceNow.Password)

	serviceNowResult, err := c.GetServiceNowResult("cmdb_ci_business_app")
	if err != nil {
		println(err)
	}

	for _, v := range serviceNowResult.Result {
		if bo, ok := v.BusinessOwner.(map[string]interface{}); ok {
			if rIct, ok := v.ResponsabileIct.(map[string]interface{}); ok {
				json.NewEncoder(out).Encode(fmt.Sprintf("serial_number: %s, business_owner_name: %v, responsabile_ict: %v, cost_center: %s \n", v.SerialNumber, bo["display_value"], rIct["display_value"], v.CostCenter))
			}
		}
	}
}
