package config

var Conf Config

type Config struct {
	ServiceNow  ServiceNow
	OracleCloud OracleCloud
}
type ServiceNow struct {
	URL      string
	Username string
	Password string
}

type OracleCloud struct {
	User        string
	Fingerprint string
	Tenancy     string
	Region      string
	Key         string
}
