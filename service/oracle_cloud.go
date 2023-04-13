package service

import (
	"context"
	"fmt"

	"github.com/ercole-io/tico/config"
	"github.com/oracle/oci-go-sdk/filestorage"
	"github.com/oracle/oci-go-sdk/loadbalancer"
	"github.com/oracle/oci-go-sdk/resourcesearch"
	"github.com/oracle/oci-go-sdk/v65/apigateway"
	"github.com/oracle/oci-go-sdk/v65/bastion"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/oracle/oci-go-sdk/v65/disasterrecovery"
	"github.com/oracle/oci-go-sdk/v65/example/helpers"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/oracle/oci-go-sdk/v65/nosql"
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

func BulkEditTags(resource resourcesearch.ResourceSummary, businessOwner string) (*identity.BulkEditTagsResponse, error) {
	client, err := identity.NewIdentityClientWithConfigurationProvider(getConfigurationProvider())
	if err != nil {
		return nil, err
	}

	req := identity.BulkEditTagsRequest{
		BulkEditTagsDetails: identity.BulkEditTagsDetails{
			BulkEditOperations: []identity.BulkEditOperationDetails{
				{
					DefinedTags:   map[string]map[string]interface{}{config.Conf.OracleCloud.OciTag.NamespaceName: {config.Conf.OracleCloud.OciTag.Name: businessOwner}},
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

func UpdateBucket(resource resourcesearch.ResourceSummary, businessOwner string) (*objectstorage.UpdateBucketResponse, error) {
	client, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(getConfigurationProvider())
	if err != nil {
		return nil, err
	}

	resource.DefinedTags[config.Conf.OracleCloud.OciTag.NamespaceName][config.Conf.OracleCloud.OciTag.Name] = businessOwner

	req := objectstorage.UpdateBucketRequest{
		BucketName:    resource.DisplayName,
		NamespaceName: common.String(config.Conf.OracleCloud.OciObjectStorage.NamespaceName),
		UpdateBucketDetails: objectstorage.UpdateBucketDetails{
			DefinedTags: resource.DefinedTags,
		}}

	resp, err := client.UpdateBucket(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func UpdateDrProtectionGroup(resource resourcesearch.ResourceSummary, businessOwner string) (*disasterrecovery.UpdateDrProtectionGroupResponse, error) {
	client, err := disasterrecovery.NewDisasterRecoveryClientWithConfigurationProvider(getConfigurationProvider())
	if err != nil {
		return nil, err
	}

	resource.DefinedTags[config.Conf.OracleCloud.OciTag.NamespaceName][config.Conf.OracleCloud.OciTag.Name] = businessOwner

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

func UpdateBastion(resource resourcesearch.ResourceSummary, businessOwner string) (*bastion.UpdateBastionResponse, error) {
	client, err := bastion.NewBastionClientWithConfigurationProvider(getConfigurationProvider())
	if err != nil {
		return nil, err
	}

	resource.DefinedTags[config.Conf.OracleCloud.OciTag.NamespaceName][config.Conf.OracleCloud.OciTag.Name] = businessOwner

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

func UpdateDrgRouteTable(resource resourcesearch.ResourceSummary, businessOwner string) (*core.UpdateDrgRouteTableResponse, error) {
	client, err := core.NewVirtualNetworkClientWithConfigurationProvider(getConfigurationProvider())
	if err != nil {
		return nil, err
	}

	resource.DefinedTags[config.Conf.OracleCloud.OciTag.NamespaceName][config.Conf.OracleCloud.OciTag.Name] = businessOwner

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

func UpdateLoadBalancer(resource resourcesearch.ResourceSummary, businessOwner string) (*loadbalancer.UpdateLoadBalancerResponse, error) {
	client, err := loadbalancer.NewLoadBalancerClientWithConfigurationProvider(getConfigurationProvider())
	if err != nil {
		return nil, err
	}

	resource.DefinedTags[config.Conf.OracleCloud.OciTag.NamespaceName][config.Conf.OracleCloud.OciTag.Name] = businessOwner

	req := loadbalancer.UpdateLoadBalancerRequest{
		LoadBalancerId: resource.Identifier,
		UpdateLoadBalancerDetails: loadbalancer.UpdateLoadBalancerDetails{
			DefinedTags: resource.DefinedTags,
		},
	}

	resp, err := client.UpdateLoadBalancer(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func UpdatePublicIpPool(resource resourcesearch.ResourceSummary, businessOwner string) (*core.UpdatePublicIpPoolResponse, error) {
	client, err := core.NewVirtualNetworkClientWithConfigurationProvider(getConfigurationProvider())
	if err != nil {
		return nil, err
	}

	resource.DefinedTags[config.Conf.OracleCloud.OciTag.NamespaceName][config.Conf.OracleCloud.OciTag.Name] = businessOwner

	req := core.UpdatePublicIpPoolRequest{
		PublicIpPoolId: resource.Identifier,
		UpdatePublicIpPoolDetails: core.UpdatePublicIpPoolDetails{
			DefinedTags: resource.DefinedTags,
		},
	}

	resp, err := client.UpdatePublicIpPool(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func UpdateFileSystem(resource resourcesearch.ResourceSummary, businessOwner string) (*filestorage.UpdateFileSystemResponse, error) {
	client, err := filestorage.NewFileStorageClientWithConfigurationProvider(getConfigurationProvider())
	if err != nil {
		return nil, err
	}

	resource.DefinedTags[config.Conf.OracleCloud.OciTag.NamespaceName][config.Conf.OracleCloud.OciTag.Name] = businessOwner

	req := filestorage.UpdateFileSystemRequest{
		FileSystemId: resource.Identifier,
		UpdateFileSystemDetails: filestorage.UpdateFileSystemDetails{
			DefinedTags: resource.DefinedTags,
		}}

	resp, err := client.UpdateFileSystem(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func UpdateApi(resource resourcesearch.ResourceSummary, businessOwner string) (*apigateway.UpdateApiResponse, error) {
	client, err := apigateway.NewApiGatewayClientWithConfigurationProvider(common.DefaultConfigProvider())
	if err != nil {
		return nil, err
	}

	resource.DefinedTags[config.Conf.OracleCloud.OciTag.NamespaceName][config.Conf.OracleCloud.OciTag.Name] = businessOwner

	req := apigateway.UpdateApiRequest{
		ApiId: resource.Identifier,
		UpdateApiDetails: apigateway.UpdateApiDetails{
			DefinedTags: resource.DefinedTags,
		}}

	resp, err := client.UpdateApi(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func UpdateComputeCapacityReservation(resource resourcesearch.ResourceSummary, businessOwner string) (*core.UpdateComputeCapacityReservationResponse, error) {
	client, err := core.NewComputeClientWithConfigurationProvider(getConfigurationProvider())
	if err != nil {
		return nil, err
	}

	resource.DefinedTags[config.Conf.OracleCloud.OciTag.NamespaceName][config.Conf.OracleCloud.OciTag.Name] = businessOwner

	req := core.UpdateComputeCapacityReservationRequest{
		CapacityReservationId: resource.Identifier,
		UpdateComputeCapacityReservationDetails: core.UpdateComputeCapacityReservationDetails{
			DefinedTags: resource.DefinedTags},
	}

	resp, err := client.UpdateComputeCapacityReservation(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func UpdateDedicatedVmHost(resource resourcesearch.ResourceSummary, businessOwner string) (*core.UpdateDedicatedVmHostResponse, error) {
	client, err := core.NewComputeClientWithConfigurationProvider(getConfigurationProvider())
	if err != nil {
		return nil, err
	}

	resource.DefinedTags[config.Conf.OracleCloud.OciTag.NamespaceName][config.Conf.OracleCloud.OciTag.Name] = businessOwner

	req := core.UpdateDedicatedVmHostRequest{
		DedicatedVmHostId: resource.Identifier,
		UpdateDedicatedVmHostDetails: core.UpdateDedicatedVmHostDetails{
			DefinedTags: resource.DefinedTags}}

	resp, err := client.UpdateDedicatedVmHost(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func UpdateTable(resource resourcesearch.ResourceSummary, businessOwner string) (*nosql.UpdateTableResponse, error) {
	client, err := nosql.NewNosqlClientWithConfigurationProvider(getConfigurationProvider())
	if err != nil {
		return nil, err
	}

	resource.DefinedTags[config.Conf.OracleCloud.OciTag.NamespaceName][config.Conf.OracleCloud.OciTag.Name] = businessOwner

	req := nosql.UpdateTableRequest{
		TableNameOrId: resource.Identifier,
		UpdateTableDetails: nosql.UpdateTableDetails{
			DefinedTags: resource.DefinedTags,
		}}

	resp, err := client.UpdateTable(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func UpdateInstance(resource resourcesearch.ResourceSummary, businessOwner string) (*core.UpdateInstanceResponse, error) {
	client, err := core.NewComputeClientWithConfigurationProvider(getConfigurationProvider())
	if err != nil {
		return nil, err
	}

	resource.DefinedTags[config.Conf.OracleCloud.OciTag.NamespaceName][config.Conf.OracleCloud.OciTag.Name] = businessOwner

	req := core.UpdateInstanceRequest{
		InstanceId: resource.Identifier,
		UpdateInstanceDetails: core.UpdateInstanceDetails{
			DefinedTags: resource.DefinedTags,
		}}

	resp, err := client.UpdateInstance(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
