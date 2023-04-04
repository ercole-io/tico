package service

import (
	"context"
	"fmt"

	"github.com/ercole-io/tico/config"
	"github.com/ercole-io/tico/model"
	"github.com/oracle/oci-go-sdk/resourcesearch"
	"github.com/oracle/oci-go-sdk/v65/bastion"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/oracle/oci-go-sdk/v65/disasterrecovery"
	"github.com/oracle/oci-go-sdk/v65/example/helpers"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/oracle/oci-go-sdk/v65/objectstorage"
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

func BulkEditTags(resource resourcesearch.ResourceSummary, snObj model.ServiceNowObj) (*identity.BulkEditTagsResponse, error) {
	client, err := identity.NewIdentityClientWithConfigurationProvider(getConfigurationProvider())
	if err != nil {
		return nil, err
	}

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
	if err != nil {
		return nil, err
	}

	return &resp, nil
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

func UpdateBucket(resource resourcesearch.ResourceSummary, snObj model.ServiceNowObj) (*objectstorage.UpdateBucketResponse, error) {
	client, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(getConfigurationProvider())
	if err != nil {
		return nil, err
	}

	businessOwner, ok := snObj.BusinessOwner.(map[string]interface{})
	if !ok {
		panic(ok)
	}

	req := objectstorage.UpdateBucketRequest{
		BucketName:    resource.DisplayName,
		NamespaceName: common.String(config.Conf.OracleCloud.OciObjectStorage.NamespaceName),
		UpdateBucketDetails: objectstorage.UpdateBucketDetails{
			DefinedTags: map[string]map[string]interface{}{config.Conf.OracleCloud.OciTag.NamespaceName: {config.Conf.OracleCloud.OciTag.Name: businessOwner["display_value"].(string)}},
		}}

	resp, err := client.UpdateBucket(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func UpdateDrProtectionGroup(resource resourcesearch.ResourceSummary, snObj model.ServiceNowObj) (*disasterrecovery.UpdateDrProtectionGroupResponse, error) {
	client, err := disasterrecovery.NewDisasterRecoveryClientWithConfigurationProvider(getConfigurationProvider())
	if err != nil {
		return nil, err
	}

	businessOwner, ok := snObj.BusinessOwner.(map[string]interface{})
	if !ok {
		panic(ok)
	}

	resource.DefinedTags[config.Conf.OracleCloud.OciTag.NamespaceName][config.Conf.OracleCloud.OciTag.Name] = businessOwner["display_value"].(string)

	req := disasterrecovery.UpdateDrProtectionGroupRequest{
		DrProtectionGroupId: common.String(config.Conf.OracleCloud.OciDrProtectionGroup.DrProtectionGroupId),
		UpdateDrProtectionGroupDetails: disasterrecovery.UpdateDrProtectionGroupDetails{
			DefinedTags: resource.DefinedTags,
		}}

	resp, err := client.UpdateDrProtectionGroup(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func UpdateBastion(resource resourcesearch.ResourceSummary, snObj model.ServiceNowObj) (*bastion.UpdateBastionResponse, error) {
	client, err := bastion.NewBastionClientWithConfigurationProvider(getConfigurationProvider())
	if err != nil {
		return nil, err
	}

	businessOwner, ok := snObj.BusinessOwner.(map[string]interface{})
	if !ok {
		panic(ok)
	}

	resource.DefinedTags[config.Conf.OracleCloud.OciTag.NamespaceName][config.Conf.OracleCloud.OciTag.Name] = businessOwner["display_value"].(string)

	req := bastion.UpdateBastionRequest{
		BastionId: resource.Identifier,
		UpdateBastionDetails: bastion.UpdateBastionDetails{
			DefinedTags: resource.DefinedTags,
		}}

	resp, err := client.UpdateBastion(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func UpdateDrgRouteTable(resource resourcesearch.ResourceSummary, snObj model.ServiceNowObj) (*core.UpdateDrgRouteTableResponse, error) {
	client, err := core.NewVirtualNetworkClientWithConfigurationProvider(getConfigurationProvider())
	if err != nil {
		return nil, err
	}

	businessOwner, ok := snObj.BusinessOwner.(map[string]interface{})
	if !ok {
		panic(ok)
	}

	resource.DefinedTags[config.Conf.OracleCloud.OciTag.NamespaceName][config.Conf.OracleCloud.OciTag.Name] = businessOwner["display_value"].(string)

	req := core.UpdateDrgRouteTableRequest{
		UpdateDrgRouteTableDetails: core.UpdateDrgRouteTableDetails{
			DefinedTags: resource.DefinedTags,
		},
		DrgRouteTableId: resource.Identifier,
	}

	resp, err := client.UpdateDrgRouteTable(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
