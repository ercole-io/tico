package model

import "github.com/oracle/oci-go-sdk/resourcesearch"

type UpdateOp struct {
	Resource       resourcesearch.ResourceSummary
	BusinnessOwner string
}
