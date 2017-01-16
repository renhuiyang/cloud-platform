package kong

import (
	"bytes"
	"errors"
	"net/http"
	"service-cloud/utils/kong/models"
	"strconv"

	"encoding/json"

	"github.com/astaxie/beego/httplib"
)

// AddAPI Add API
// name -- optional	The API name. If none is specified, will default to the request_host or request_path.
// request_host -- semi-optional	The public DNS address that points to your API. For example, mockbin.com. At least request_host or request_path or both should be specified.
// request_path -- emi-optional	The public path that points to your API. For example, /someservice. At least request_host or request_path or both should be specified.
// strip_request_path -- optional	Strip the request_path value before proxying the request to the final API. For example a request made to /someservice/hello will be resolved to upstream_url/hello. By default is false.
// preserve_host -- optional	Preserves the original Host header sent by the client, instead of replacing it with the hostname of the upstream_url. By default is false.
// upstream_url	The base target URL that points to your API server, this URL will be used for proxying requests. For example, https://mockbin.com.
func AddAPI(api *models.API) (*models.API, error) {
	// POST /apis/
	req := httplib.Post(kongAdminURL + `/apis/`)
	req.Param("name", api.Name)
	req.Param("request_host", api.RequestHost)
	req.Param("request_path", api.RequestPath)
	req.Param("strip_request_path", strconv.FormatBool(api.StripRequestPath))
	req.Param("preserve_host", strconv.FormatBool(api.PreserveHost))
	req.Param("upstream_url", api.UpstreamURL)

	var retAPI models.API
	err := req.ToJSON(&retAPI)
	if err != nil {
		return nil, err
	}
	return &retAPI, nil
}

// GetAPI Retrieve API
// nameOrID -- (required)The unique identifier or the name of the API to retrieve.
func GetAPI(nameOrID string) (*models.API, error) {
	//GET /apis/{name or id}
	if len(nameOrID) == 0 {
		return nil, errors.New("The unique identifier or the name of the API can not be null")
	}
	req := httplib.Get(kongAdminURL + `/apis/` + nameOrID)

	var retAPI models.API
	err := req.ToJSON(&retAPI)
	if err != nil {
		return nil, err
	}
	return &retAPI, nil
}

//ListAPIs List APIs
func ListAPIs(size int, offset string) (*models.APIList, error) {
	//GET /apis/
	return nil, nil
}

// UpdateAPI Update API
// nameOrID (required) The original unique identifier or the name of the API to update
// name -- (optional)The API name. If none is specified, will default to the request_host or request_path.
// request_host -- (semi-optional)The public DNS address that points to your API. For example, mockbin.com. At least request_host or request_path or both should be specified.
// request_path -- (semi-optional)	The public path that points to your API. For example, /someservice. At least request_host or request_path or both should be specified.
// strip_request_path -- (optional)	Strip the request_path value before proxying the request to the final API. For example a request made to /someservice/hello will be resolved to upstream_url/hello. By default is false.
// preserve_host -- (optional)	Preserves the original Host header sent by the client, instead of replacing it with the hostname of the upstream_url. By default is false.
// upstream_url	The base target URL that points to your API server, this URL will be used for proxying requests. For example, https://mockbin.com.
func UpdateAPI(nameOrID string, api *models.API) (*models.API, error) {
	//PATCH /apis/{name or id}
	if len(nameOrID) == 0 {
		return nil, errors.New("The unique identifier or the name of the API can not be null")
	}

	jsonStr, err := json.Marshal(api)
	url := kongAdminURL + `/apis/` + nameOrID
	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonStr))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var retAPI models.API
	err = json.NewDecoder(resp.Body).Decode(&retAPI)
	if err != nil {
		return nil, err
	}

	return &retAPI, nil
}

//DeleteAPI Delete API
func DeleteAPI(nameOrID string) error {
	//DELETE /apis/{name or id}

	if len(nameOrID) == 0 {
		return errors.New("The unique identifier or the name of the API can not be null")
	}
	req := httplib.Delete(kongAdminURL + `/apis/` + nameOrID)

	_, err := req.Response()
	if err != nil {
		return err
	}
	return nil
}
