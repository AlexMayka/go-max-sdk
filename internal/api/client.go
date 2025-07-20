package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlexMayka/go-max-sdk/internal/config"
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
		return nil, nil, nil, fmt.Errorf("неправильный тип запроса")
	}

	pathParams := make(map[string]string)
	queryParams := make(map[string]string)
	jsonBody := make(map[string]interface{})

	v := reflect.ValueOf(req)
	t := reflect.TypeOf(req)

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if pathTag := field.Tag.Get("path"); pathTag != "" {
			pathParams[pathTag] = getStringValue(value)
			continue
		}

		if queryTag := field.Tag.Get("query"); queryTag != "" {
			tag := strings.Split(queryTag, ",")
			hasOmitEmpty := len(tag) == 2 && tag[1] == "omitempty"

			if hasOmitEmpty && shouldOmit(value) {
				continue
			}

			queryParams[tag[0]] = getStringValue(value)
			continue
		}

		if jsonTag := field.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
			tag := strings.Split(jsonTag, ",")
			hasOmitEmpty := len(tag) == 2 && tag[1] == "omitempty"

			if hasOmitEmpty && shouldOmit(value) {
				continue
			}

			jsonBody[tag[0]] = getInterfaceValue(value)
		}
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

	values := address.Query()

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

func shouldOmit(value reflect.Value) bool {
	if value.Kind() == reflect.Ptr && value.IsNil() {
		return true
	}
	if value.Kind() != reflect.Ptr && value.IsZero() {
		return true
	}
	return false
}

func getInterfaceValue(value reflect.Value) interface{} {
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil
		}
		return value.Elem().Interface()
	}
	return value.Interface()
}

func getStringValue(value reflect.Value) string {
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return ""
		}
		return fmt.Sprintf("%v", value.Elem().Interface())
	}
	return fmt.Sprintf("%v", value.Interface())
}
