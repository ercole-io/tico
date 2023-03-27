package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/ercole-io/tico/api"
	"github.com/ercole-io/tico/config"
	"github.com/ercole-io/tico/service"
)

func Handler(ctx context.Context, in io.Reader, out io.Writer) {
	resp, err := service.CreateTag(config.Conf.OracleCloud.OciTag.NamespaceId, config.Conf.OracleCloud.OciTag.Name, config.Conf.OracleCloud.OciTag.Description)
	if resp.RawResponse.StatusCode == 409 {
		log.Println(err)
	}

	c := api.New(config.Conf.ServiceNow.URL, config.Conf.ServiceNow.Username, config.Conf.ServiceNow.Password)

	snList, err := c.GetServiceNowResult(config.Conf.ServiceNow.TableName)
	if err != nil {
		log.Println(err)
	}

	log.Printf("ServiceNow items fetched: %v", len(snList.Result))

	ocList := service.SearchResources(config.Conf.OracleCloud.Match.Element)

	log.Printf("Oracle Cloud Tags items fetched: %v", len(ocList))

	updateCounter := 0
	for _, oc := range ocList {
		for _, tag := range oc.DefinedTags {
			if v, ok := tag[config.Conf.OracleCloud.Match.Element]; ok {
				for _, sn := range snList.Result {
					if sn.SerialNumber == v {
						log.Printf("ServiceNow item SerialNumber: %s  -  Oracle Cloud resource to updated: %s", sn.SerialNumber, *oc.Identifier)
						resp := service.BulkEditTags(oc, sn)
						log.Printf("Edit tag response: %s", resp.RawResponse.Status)
						updateCounter++
					}
				}
			}
		}
	}

	outCounter := fmt.Sprintf("Oracle Cloud updated tags:%v", updateCounter)

	log.Print(outCounter)

	json.NewEncoder(out).Encode(outCounter)
}
