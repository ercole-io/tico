package service

import (
	"context"
	"fmt"

	"github.com/ercole-io/tico/config"
	"github.com/ercole-io/tico/model"
	"github.com/oracle/oci-go-sdk/resourcesearch"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/example/helpers"
	"github.com/oracle/oci-go-sdk/v65/identity"
)

func getConfigurationProvider() common.ConfigurationProvider {
	return common.NewRawConfigurationProvider(
		config.Conf.OracleCloud.Tenancy,
		config.Conf.OracleCloud.User,
		config.Conf.OracleCloud.Region,
		config.Conf.OracleCloud.Fingerprint,
		config.Conf.OracleCloud.Key, nil)
}

func SearchResources(definedTagKey string) []resourcesearch.ResourceSummary {
	client, err := resourcesearch.NewResourceSearchClientWithConfigurationProvider(getConfigurationProvider())
	helpers.FatalIfError(err)

	req := resourcesearch.SearchResourcesRequest{SearchDetails: resourcesearch.StructuredSearchDetails{
		Query: common.String(fmt.Sprintf("query all resources where definedTags.key = '%s'", definedTagKey))},
	}

	resp, err := client.SearchResources(context.Background(), req)
	helpers.FatalIfError(err)

	return resp.ResourceSummaryCollection.Items
}

func BulkEditTags(resource resourcesearch.ResourceSummary, snObj model.ServiceNowObj) identity.BulkEditTagsResponse {
	client, err := identity.NewIdentityClientWithConfigurationProvider(getConfigurationProvider())
	helpers.FatalIfError(err)

	businessOwner, ok := snObj.BusinessOwner.(map[string]interface{})
	if !ok {
		panic(ok)
	}

	req := identity.BulkEditTagsRequest{BulkEditTagsDetails: identity.BulkEditTagsDetails{
		BulkEditOperations: []identity.BulkEditOperationDetails{
			{
				DefinedTags:   map[string]map[string]interface{}{config.Conf.OracleCloud.OciTag.NamespaceName: {config.Conf.OracleCloud.OciTag.Name: businessOwner["display_value"].(string)}},
				OperationType: identity.BulkEditOperationDetailsOperationTypeAddOrSet,
			},
		},
		CompartmentId: resource.CompartmentId,
		Resources: []identity.BulkEditResource{
			{
				Id:           resource.Identifier,
				ResourceType: resource.ResourceType,
			},
		},
	},
	}

	resp, err := client.BulkEditTags(context.Background(), req)
	helpers.FatalIfError(err)

	return resp
}

func CreateTag(tagNamespaceId string, name string, description string) (identity.CreateTagResponse, error) {
	client, err := identity.NewIdentityClientWithConfigurationProvider(getConfigurationProvider())
	helpers.FatalIfError(err)

	req := identity.CreateTagRequest{
		CreateTagDetails: identity.CreateTagDetails{
			Description: common.String(description),
			Validator:   identity.DefaultTagDefinitionValidator{},
			Name:        common.String(name)},
		TagNamespaceId: common.String(tagNamespaceId)}

	return client.CreateTag(context.Background(), req)
}
