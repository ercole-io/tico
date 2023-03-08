package model

type ServiceNowResult struct {
	Result []ServiceNowObj `json:"result"`
}

type ServiceNowObj struct {
	SerialNumber    string      `json:"serial_number"`
	BusinessOwner   interface{} `json:"u_uo_business_owner,omitempty"`
	ResponsabileIct interface{} `json:"u_uo_responsabile_ict,omitempty"`
	CostCenter      string      `json:"cost_center"`
}
