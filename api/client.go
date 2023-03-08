package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ercole-io/tico/model"
)

type Client struct {
	url      string
	user     string
	password string
	client   *http.Client
}

func New(url string, user string, password string) *Client {
	return &Client{
		url:      url,
		user:     user,
		password: password,
		client:   http.DefaultClient,
	}
}

func (c *Client) GetServiceNowResult(sysClassName string) (*model.ServiceNowResult, error) {
	queryParam := url.Values{}
	queryParam.Add("sysparm_query", fmt.Sprintf("sys_class_name=%s", sysClassName))
	queryParam.Add("sysparm_display_value", "true")

	header := http.Header{}
	header.Add("Authorization", c.basicAuth(c.user, c.password))

	resp, err := c.doRequest(http.MethodGet, "api/now/table/cmdb_ci_business_app", queryParam, header, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &model.ServiceNowResult{}

	if err := json.Unmarshal(body, result); err != nil { // Parse []byte to go struct pointer
		fmt.Printf("Can not unmarshal JSON, %s", err)
	}

	return result, nil
}

func (c *Client) basicAuth(username, password string) string {
	auth := username + ":" + password
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(auth)))
}

func (c *Client) doRequest(method, path string, query url.Values, header http.Header, ibody io.Reader) (*http.Response, error) {
	u, err := url.Parse(c.url + "/" + path)
	if err != nil {
		return nil, err
	}
	u.RawQuery = query.Encode()

	req, err := http.NewRequest(method, u.String(), ibody)
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		req.Header[k] = v
	}

	res, err := c.client.Do(req)

	return res, err
}
