package controller

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/ercole-io/tico/api"
	"github.com/ercole-io/tico/config"
	"github.com/ercole-io/tico/model"
	"github.com/ercole-io/tico/service"
	"github.com/oracle/oci-go-sdk/resourcesearch"
)

func Handler(ctx context.Context, in io.Reader, out io.Writer) {
	resp, _ := service.CreateTag(config.Conf.OracleCloud.OciTag.NamespaceId, config.Conf.OracleCloud.OciTag.Name, config.Conf.OracleCloud.OciTag.Description)
	if resp.RawResponse.StatusCode != 409 {
		log.Printf("created new tag %s: %s", config.Conf.OracleCloud.OciTag.Name, resp.RawResponse.Status)
	}

	snList := make(chan *model.ServiceNowResult)
	ocList := make(chan []resourcesearch.ResourceSummary)

	c := api.New(config.Conf.ServiceNow.URL, config.Conf.ServiceNow.Username, config.Conf.ServiceNow.Password)

	go func(c *api.Client) {
		res, err := c.GetServiceNowResult(config.Conf.ServiceNow.TableName)
		if err != nil {
			log.Println(err)
		}

		snList <- res
	}(c)

	go func() {
		ocList <- service.SearchResources(config.Conf.OracleCloud.Match.Element)
	}()

	log.Println(updateItems(<-snList, <-ocList))
}
func printItemsSize(s *model.ServiceNowResult, o []resourcesearch.ResourceSummary) {
	log.Printf("ServiceNow items fetched: %v", len(s.Result))
	log.Printf("Oracle Cloud Tags items fetched: %v", len(o))
}

func updateItems(snList *model.ServiceNowResult, ocList []resourcesearch.ResourceSummary) string {
	printItemsSize(snList, ocList)

	updateCounter := 0
	for _, oc := range ocList {
	DefinedTag:
		for _, tag := range oc.DefinedTags {
			if v, ok := tag[config.Conf.OracleCloud.Match.Element]; ok {
				for _, sn := range snList.Result {
					if sn.SerialNumber == v {
						log.Printf("ServiceNow item SerialNumber: %s  -  Oracle Cloud resource to updated: %s", sn.SerialNumber, *oc.DisplayName)
						switch *oc.ResourceType {
						case "Bucket":
							resp, err := service.UpdateBucket(oc, sn)
							if err != nil {
								log.Print(err)
							} else {
								log.Printf("Edit tag response: %s", resp.RawResponse.Status)
								updateCounter++
							}
							break DefinedTag
						case "DrProtectionGroup":
							resp, err := service.UpdateDrProtectionGroup(oc, sn)
							if err != nil {
								log.Print(err)
							} else {
								log.Printf("Edit tag response: %s", resp.RawResponse.Status)
								updateCounter++
							}
							break DefinedTag
						case "Bastion":
							resp, err := service.UpdateBastion(oc, sn)
							if err != nil {
								log.Print(err)
							} else {
								log.Printf("Edit tag response: %s", resp.RawResponse.Status)
								updateCounter++
							}
							break DefinedTag
						case "DrgRouteTable":
							resp, err := service.UpdateDrgRouteTable(oc, sn)
							if err != nil {
								log.Print(err)
							} else {
								log.Printf("Edit tag response: %s", resp.RawResponse.Status)
								updateCounter++
							}
							break DefinedTag
						default:
							resp, err := service.BulkEditTags(oc, sn)
							if err != nil {
								log.Print(err)
							} else {
								log.Printf("Edit tag response: %s", resp.RawResponse.Status)
								updateCounter++
							}
							break DefinedTag
						}

					}
				}
			}
		}
	}

	return fmt.Sprintf("Oracle Cloud updated tags:%v", updateCounter)
}
