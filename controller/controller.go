package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ercole-io/tico/api"
	"github.com/ercole-io/tico/config"
	"github.com/ercole-io/tico/service"
)

func Handler(ctx context.Context, in io.Reader, out io.Writer) {
	resp, err := service.CreateTag(config.Conf.OracleCloud.OciTag.NamespaceId, config.Conf.OracleCloud.OciTag.Name, config.Conf.OracleCloud.OciTag.Description)
	if resp.RawResponse.StatusCode == http.StatusConflict {
		fmt.Sprintln(err)
	}

	c := api.New(config.Conf.ServiceNow.URL, config.Conf.ServiceNow.Username, config.Conf.ServiceNow.Password)

	snList, err := c.GetServiceNowResult(config.Conf.ServiceNow.TableName)
	if err != nil {
		println(err)
	}

	ocList := service.SearchResources(config.Conf.OracleCloud.Match.Element)

	for _, oc := range ocList {
		for _, tag := range oc.DefinedTags {
			if v, ok := tag[config.Conf.OracleCloud.Match.Element]; ok {
				for _, sn := range snList.Result {
					if sn.SerialNumber == v {
						resp := service.BulkEditTags(oc, sn)
						println(resp.RawResponse.StatusCode)
					}
				}

			}
		}

	}

	json.NewEncoder(out).Encode("end")

}
