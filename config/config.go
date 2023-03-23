package config

var Conf Config

type Config struct {
	ServiceNow  ServiceNow
	OracleCloud OracleCloud
}
type ServiceNow struct {
	URL       string
	Username  string
	Password  string
	TableName string
	Match     Match
}

type OracleCloud struct {
	User        string
	Fingerprint string
	Tenancy     string
	Region      string
	Key         string
	OciTag      OciTag
	Match       Match
}

type OciTag struct {
	NamespaceId   string
	NamespaceName string
	Name          string
	Description   string
}

type Match struct {
	Element  string
	Elements []string
}
