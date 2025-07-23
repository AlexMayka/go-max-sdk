package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlexMayka/go-max-sdk/internal/config"
	"github.com/AlexMayka/go-max-sdk/internal/utils"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type Client struct {
	token      string
	httpClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		token:      token,
		httpClient: &http.Client{},
	}
}

func (c *Client) parseRequest(req interface{}, cfg *config.EndpointConfig) (map[string]string, map[string]string, interface{}, error) {
	reqType := reflect.TypeOf(req)
	if reqType != cfg.RequestModel {
		return nil, nil, nil, fmt.Errorf("invalid request type")
	}

	pathParams := make(map[string]string)
	queryParams := make(map[string]string)
	jsonBody := make(map[string]interface{})

	v := reflect.ValueOf(req)
	t := reflect.TypeOf(req)

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		c.parsePathParams(field, value, pathParams)
		c.parseQueryParams(field, value, queryParams)
		c.parseJSONBody(field, value, jsonBody)
	}

	return pathParams, queryParams, jsonBody, nil
}

func (c *Client) buildURL(path string, pathParams, queryParams map[string]string) string {
	address := url.URL{}

	address.Scheme = config.Scheme
	address.Host = config.Host

	for k, v := range pathParams {
		path = strings.Replace(path, fmt.Sprintf("{%s}", k), v, -1)
	}
	address.Path = path

	values := url.Values{}

	if c.token != "" {
		values.Add("access_token", c.token)
	}

	for k, v := range queryParams {
		values.Add(k, v)
	}
	address.RawQuery = values.Encode()

	return address.String()
}

func sendRequest(ctx context.Context, jsonBody interface{}, cfg *config.EndpointConfig, addr string) (interface{}, error) {
	var body []byte
	var err error

	if jsonBody != nil {
		body, err = json.Marshal(jsonBody)
		if err != nil {
			return nil, err
		}
	}

	request, err := http.NewRequestWithContext(ctx, string(cfg.Method), addr, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	if len(body) > 0 {
		request.Header.Set("Content-Type", string(cfg.ContentType))
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d - %s", response.StatusCode, string(responseBody))
	}

	responsePtr := reflect.New(cfg.ResponseModel)
	responseValue := responsePtr.Interface()

	err = json.Unmarshal(responseBody, responseValue)
	if err != nil {
		return nil, err
	}

	return responseValue, nil
}

func (c *Client) Call(ctx context.Context, endpoint config.Endpoint, req interface{}) (interface{}, error) {
	cfg := config.EndpointConfigs[endpoint]

	pathParams, queryParams, jsonBody, err := c.parseRequest(req, cfg)
	if err != nil {
		return nil, err
	}

	addr := c.buildURL(cfg.Path, pathParams, queryParams)

	answer, err := sendRequest(ctx, jsonBody, cfg, addr)
	if err != nil {
		return nil, err
	}

	return answer, nil
}

func (c *Client) parsePathParams(field reflect.StructField, value reflect.Value, pathParams map[string]string) {
	pathTag := field.Tag.Get("path")
	if pathTag != "" {
		pathParams[pathTag] = utils.GetStringValue(value)
	}
}

func (c *Client) parseQueryParams(field reflect.StructField, value reflect.Value, queryParams map[string]string) {
	queryTag := field.Tag.Get("query")
	if queryTag == "" {
		return
	}

	tag := strings.Split(queryTag, ",")
	hasOmitEmpty := len(tag) == 2 && tag[1] == "omitempty"

	if hasOmitEmpty && utils.ShouldOmit(value) {
		return
	}

	queryParams[tag[0]] = utils.GetStringValue(value)
}

func (c *Client) parseJSONBody(field reflect.StructField, value reflect.Value, jsonBody map[string]interface{}) {
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" || jsonTag == "-" {
		return
	}

	tag := strings.Split(jsonTag, ",")
	hasOmitEmpty := len(tag) == 2 && tag[1] == "omitempty"

	if hasOmitEmpty && utils.ShouldOmit(value) {
		return
	}

	jsonBody[tag[0]] = utils.GetInterfaceValue(value)
}
