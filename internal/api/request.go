package api

import (
	"fmt"
	"github.com/AlexMayka/go-max-sdk/internal/utils"
	"net/url"
	"reflect"
	"strings"
)

// buildURL constructs the complete API URL by combining scheme, host, path parameters, and query parameters.
// It replaces path placeholders (e.g., {chatId}) with actual values and adds authentication token.
// Parameters:
//   - path: URL path template with {param} placeholders
//   - pathParams: map of path parameter names to values
//   - queryParams: map of query parameter names to values
// Returns: Complete URL string ready for HTTP request
func (c *Client) buildURL(path string, pathParams, queryParams map[string]string) string {
	address := url.URL{}

	address.Scheme = Scheme
	address.Host = Host

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

// parseRequest analyzes a request struct and extracts path parameters, query parameters, and JSON body.
// It uses reflection to inspect struct fields and their tags to determine how to process each field:
//   - "path" tag: field value becomes a path parameter
//   - "query" tag: field value becomes a query parameter
//   - "json" tag: field value becomes part of JSON request body
// Parameters:
//   - req: request structure to parse
//   - cfg: endpoint configuration containing expected request type
// Returns:
//   - pathParams: map of path parameter names to values
//   - queryParams: map of query parameter names to values
//   - jsonBody: interface{} containing the JSON request body
//   - error: validation error if request type doesn't match expected type
func (c *Client) parseRequest(req interface{}, cfg *EndpointConfig) (map[string]string, map[string]string, interface{}, error) {
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

// parsePathParams extracts path parameters from struct fields marked with "path" tag.
// Path parameters are used to replace placeholders in URL paths (e.g., /chats/{chatId}).
// Parameters:
//   - field: struct field metadata containing tags
//   - value: actual field value
//   - pathParams: map to store extracted path parameters
func (c *Client) parsePathParams(field reflect.StructField, value reflect.Value, pathParams map[string]string) {
	pathTag := field.Tag.Get("path")
	if pathTag != "" {
		pathParams[pathTag] = utils.GetStringValue(value)
	}
}

// parseQueryParams extracts query parameters from struct fields marked with "query" tag.
// Supports "omitempty" option to skip fields with zero values.
// Query parameters are appended to the URL after the "?" symbol.
// Parameters:
//   - field: struct field metadata containing tags
//   - value: actual field value
//   - queryParams: map to store extracted query parameters
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

// parseJSONBody extracts JSON body fields from struct fields marked with "json" tag.
// Supports "omitempty" option to skip fields with zero values.
// Fields marked with json:"-" are ignored completely.
// Parameters:
//   - field: struct field metadata containing tags
//   - value: actual field value
//   - jsonBody: map to store JSON body fields
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
