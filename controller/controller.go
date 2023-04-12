package controller

import (
	"context"
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

	u1 := make(chan int)
	u2 := make(chan int)

	firstList := <-snList
	secList := <-ocList

	go func(snHalfL *model.ServiceNowResult, ocL []resourcesearch.ResourceSummary) {
		u1 <- updateItems(snHalfL.Result[:len(snHalfL.Result)/2], ocL)
	}(firstList, secList)

	go func(snHalfL *model.ServiceNowResult, ocL []resourcesearch.ResourceSummary) {
		u2 <- updateItems(snHalfL.Result[len(snHalfL.Result)/2:len(snHalfL.Result)], ocL)
	}(firstList, secList)

	log.Printf("Oracle Cloud updated tags:%v", <-u1+<-u2)
}

func updateItems(snList []model.ServiceNowObj, ocList []resourcesearch.ResourceSummary) int {
	updateCounter := 0
	for _, oc := range ocList {
	DefinedTag:
		for _, tag := range oc.DefinedTags {
			if v, ok := tag[config.Conf.OracleCloud.Match.Element]; ok {
				for _, sn := range snList {
					if sn.SerialNumber == v {
						log.Print(*oc.ResourceType)
						log.Printf("ServiceNow item Serialnumber: %s  -  Oracle Cloud resource to updated: %s/%s", sn.SerialNumber, *oc.ResourceType, *oc.DisplayName)
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
						case "LoadBalancer":
							resp, err := service.UpdateLoadBalancer(oc, sn)
							if err != nil {
								log.Print(err)
							} else {
								log.Printf("Edit tag response: %s", resp.RawResponse.Status)
								updateCounter++
							}
							break DefinedTag
						case "PublicIpPool":
							resp, err := service.UpdatePublicIpPool(oc, sn)
							if err != nil {
								log.Print(err)
							} else {
								log.Printf("Edit tag response: %s", resp.RawResponse.Status)
								updateCounter++
							}
							break DefinedTag
						case "FileSystem":
							resp, err := service.UpdateFileSystem(oc, sn)
							if err != nil {
								log.Print(err)
							} else {
								log.Printf("Edit tag response: %s", resp.RawResponse.Status)
								updateCounter++
							}
							break DefinedTag
						case "ApiGatewayApi":
							resp, err := service.UpdateApi(oc, sn)
							if err != nil {
								log.Print(err)
							} else {
								log.Printf("Edit tag response: %s", resp.RawResponse.Status)
								updateCounter++
							}
							break DefinedTag
						case "ComputeCapacityReservation":
							resp, err := service.UpdateComputeCapacityReservation(oc, sn)
							if err != nil {
								log.Print(err)
							} else {
								log.Printf("Edit tag response: %s", resp.RawResponse.Status)
								updateCounter++
							}
							break DefinedTag
						case "DedicatedVmHost":
							resp, err := service.UpdateDedicatedVmHost(oc, sn)
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

	return updateCounter
}
