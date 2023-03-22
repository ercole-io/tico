# tico
The Tico project can sync serviceNow tags with Oracle Cloud

## Getting started
Create a oracle cloud [Application function](https://docs.oracle.com/en-us/iaas/Content/Functions/Tasks/functionscreatingfirst.htm).

In the application clone tico

`git clone https://github.com/ercole-io/tico.git`

Enter the project

`cd tico`
 
create the configuration:

```
[ServiceNow]
URL = "https://<YOUR INSTANCE>.service-now.com"
Username = ""
Password = ""

[OracleCloud]
User=""
Region=""
Key="""-----BEGIN PRIVATE KEY-----
<YOUR KEY>
-----END PRIVATE KEY-----"""
```

Deploy your application

`fn -v deploy --app <APPLICATION NAME>`

Invoke

`fn invoke <APPLICATION NAME> tico`

or using the specific invoke endpoint

```
oci raw-request --http-method POST --target-uri <INVOKE ENDPOINT> --request-body ""
```
