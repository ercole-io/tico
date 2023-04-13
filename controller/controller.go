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

	firstList := <-snList
	secList := <-ocList

	toUpdateCh := make(chan []model.UpdateOp)

	go func(snHalfL *model.ServiceNowResult, ocL []resourcesearch.ResourceSummary) {
		toUpdateCh <- getItemToUpdate(snHalfL.Result[:len(snHalfL.Result)/2], ocL)
	}(firstList, secList)

	go func(snHalfL *model.ServiceNowResult, ocL []resourcesearch.ResourceSummary) {
		toUpdateCh <- getItemToUpdate(snHalfL.Result[len(snHalfL.Result)/2:len(snHalfL.Result)], ocL)
	}(firstList, secList)

	if err := updateItems(<-toUpdateCh); err != nil {
		log.Println(err)
	}
}

func getItemToUpdate(snList []model.ServiceNowObj, ocList []resourcesearch.ResourceSummary) []model.UpdateOp {
	toUpdateList := make([]model.UpdateOp, 0, len(ocList))

	for _, oc := range ocList {
	DefinedTag:
		for _, tag := range oc.DefinedTags {
			if v, ok := tag[config.Conf.OracleCloud.Match.Element]; ok {
				for _, sn := range snList {
					if sn.SerialNumber == v {
						businessOwner, ok := sn.BusinessOwner.(map[string]interface{})
						if !ok {
							panic(ok)
						}

						toUpdateList = append(toUpdateList, model.UpdateOp{
							Resource: oc, BusinnessOwner: businessOwner["display_value"].(string)})

						break DefinedTag
					}
				}
			}
		}
	}

	return toUpdateList
}

func updateItems(toUpdateList []model.UpdateOp) error {
	for _, up := range toUpdateList {
		switch *up.Resource.ResourceType {
		case "Bucket":
			resp, err := service.UpdateBucket(up.Resource, up.BusinnessOwner)
			if err != nil {
				return err
			}

			log.Printf("Edit tag response of resource %s/%s: %s", *up.Resource.ResourceType, *up.Resource.DisplayName, resp.RawResponse.Status)
		case "DrProtectionGroup":
			resp, err := service.UpdateDrProtectionGroup(up.Resource, up.BusinnessOwner)
			if err != nil {
				return err
			}

			log.Printf("Edit tag response of resource %s/%s: %s", *up.Resource.ResourceType, *up.Resource.DisplayName, resp.RawResponse.Status)
		case "Bastion":
			resp, err := service.UpdateBastion(up.Resource, up.BusinnessOwner)
			if err != nil {
				return err
			}

			log.Printf("Edit tag response of resource %s/%s: %s", *up.Resource.ResourceType, *up.Resource.DisplayName, resp.RawResponse.Status)
		case "DrgRouteTable":
			resp, err := service.UpdateDrgRouteTable(up.Resource, up.BusinnessOwner)
			if err != nil {
				return err
			}

			log.Printf("Edit tag response of resource %s/%s: %s", *up.Resource.ResourceType, *up.Resource.DisplayName, resp.RawResponse.Status)
		case "LoadBalancer":
			resp, err := service.UpdateLoadBalancer(up.Resource, up.BusinnessOwner)
			if err != nil {
				return err
			}

			log.Printf("Edit tag response of resource %s/%s: %s", *up.Resource.ResourceType, *up.Resource.DisplayName, resp.RawResponse.Status)
		case "PublicIpPool":
			resp, err := service.UpdatePublicIpPool(up.Resource, up.BusinnessOwner)
			if err != nil {
				return err
			}

			log.Printf("Edit tag response of resource %s/%s: %s", *up.Resource.ResourceType, *up.Resource.DisplayName, resp.RawResponse.Status)
		case "FileSystem":
			resp, err := service.UpdateFileSystem(up.Resource, up.BusinnessOwner)
			if err != nil {
				return err
			}

			log.Printf("Edit tag response of resource %s/%s: %s", *up.Resource.ResourceType, *up.Resource.DisplayName, resp.RawResponse.Status)
		case "ApiGatewayApi":
			resp, err := service.UpdateApi(up.Resource, up.BusinnessOwner)
			if err != nil {
				return err
			}

			log.Printf("Edit tag response of resource %s/%s: %s", *up.Resource.ResourceType, *up.Resource.DisplayName, resp.RawResponse.Status)
		case "ComputeCapacityReservation":
			resp, err := service.UpdateComputeCapacityReservation(up.Resource, up.BusinnessOwner)
			if err != nil {
				return err
			}

			log.Printf("Edit tag response of resource %s/%s: %s", *up.Resource.ResourceType, *up.Resource.DisplayName, resp.RawResponse.Status)
		case "DedicatedVmHost":
			resp, err := service.UpdateDedicatedVmHost(up.Resource, up.BusinnessOwner)
			if err != nil {
				return err
			}

			log.Printf("Edit tag response of resource %s/%s: %s", *up.Resource.ResourceType, *up.Resource.DisplayName, resp.RawResponse.Status)
		case "NoSQLTable":
			resp, err := service.UpdateTable(up.Resource, up.BusinnessOwner)
			if err != nil {
				return err
			}

			log.Printf("Edit tag response of resource %s/%s: %s", *up.Resource.ResourceType, *up.Resource.DisplayName, resp.RawResponse.Status)
		case "Instance":
			resp, err := service.UpdateInstance(up.Resource, up.BusinnessOwner)
			if err != nil {
				return err
			}

			log.Printf("Edit tag response of resource %s/%s: %s", *up.Resource.ResourceType, *up.Resource.DisplayName, resp.RawResponse.Status)
		default:
			resp, err := service.BulkEditTags(up.Resource, up.BusinnessOwner)
			if err != nil {
				return err
			}

			log.Printf("Edit tag response of resource %s/%s: %s", *up.Resource.ResourceType, *up.Resource.DisplayName, resp.RawResponse.Status)
		}
	}

	return nil
}
