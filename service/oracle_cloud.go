package service

import (
	"context"

	"github.com/ercole-io/tico/config"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/example/helpers"
	"github.com/oracle/oci-go-sdk/resourcesearch"
)

func ExampleSearchResources() []resourcesearch.ResourceSummary {
	configurationProvider := common.NewRawConfigurationProvider(
		config.Conf.OracleCloud.Tenancy,
		config.Conf.OracleCloud.User,
		config.Conf.OracleCloud.Region,
		config.Conf.OracleCloud.Fingerprint,
		config.Conf.OracleCloud.Key, nil)

	client, err := resourcesearch.NewResourceSearchClientWithConfigurationProvider(configurationProvider)
	helpers.FatalIfError(err)

	req := resourcesearch.SearchResourcesRequest{SearchDetails: resourcesearch.StructuredSearchDetails{
		Query: common.String("query all resources")},
	}

	resp, err := client.SearchResources(context.Background(), req)
	helpers.FatalIfError(err)

	return resp.ResourceSummaryCollection.Items
}
