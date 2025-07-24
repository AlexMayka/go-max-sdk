package api

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/AlexMayka/go-max-sdk/internal/config"
)

type TestRequest struct {
	PathParam     string  `path:"id"`
	QueryParam    string  `query:"filter"`
	OptionalQuery *string `query:"optional,omitempty"`
	JSONField     string  `json:"data"`
	OptionalJSON  *string `json:"optional,omitempty"`
	IgnoredField  string  `json:"-"`
}

func TestParseRequest(t *testing.T) {
	client := NewClient("test-token").(*Client)

	cfg := &config.EndpointConfig{
		RequestModel: reflect.TypeOf((*TestRequest)(nil)).Elem(),
	}

	t.Run("ParseAllFields", func(t *testing.T) {
		optional := "optional_value"
		req := TestRequest{
			PathParam:     "123",
			QueryParam:    "active",
			OptionalQuery: &optional,
			JSONField:     "test_data",
			OptionalJSON:  &optional,
			IgnoredField:  "ignored",
		}

		pathParams, queryParams, jsonBody, err := client.parseRequest(req, cfg)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if pathParams["id"] != "123" {
			t.Errorf("Expected path param id=123, got %s", pathParams["id"])
		}

		if queryParams["filter"] != "active" {
			t.Errorf("Expected query param filter=active, got %s", queryParams["filter"])
		}
		if queryParams["optional"] != "optional_value" {
			t.Errorf("Expected query param optional=optional_value, got %s", queryParams["optional"])
		}

		jsonMap := jsonBody.(map[string]interface{})
		if jsonMap["data"] != "test_data" {
			t.Errorf("Expected json field data=test_data, got %v", jsonMap["data"])
		}
		if jsonMap["optional"] != "optional_value" {
			t.Errorf("Expected json field optional=optional_value, got %v", jsonMap["optional"])
		}

		if _, exists := jsonMap["IgnoredField"]; exists {
			t.Error("Expected ignored field to not be present in JSON body")
		}
	})

	t.Run("OmitEmptyFields", func(t *testing.T) {
		req := TestRequest{
			PathParam:  "123",
			QueryParam: "active",
			JSONField:  "test_data",
		}

		pathParams, queryParams, jsonBody, err := client.parseRequest(req, cfg)
		_ = pathParams

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if _, exists := queryParams["optional"]; exists {
			t.Error("Expected optional query param to be omitted")
		}

		jsonMap := jsonBody.(map[string]interface{})
		if _, exists := jsonMap["optional"]; exists {
			t.Error("Expected optional json field to be omitted")
		}
	})

	t.Run("WrongRequestType", func(t *testing.T) {
		type WrongRequest struct {
			Field string `json:"field"`
		}

		req := WrongRequest{Field: "test"}

		_, _, _, err := client.parseRequest(req, cfg)

		if err == nil {
			t.Error("Expected error for wrong request type")
		}
	})
}

func TestBuildURL(t *testing.T) {
	client := NewClient("test-token").(*Client)

	t.Run("WithPathAndQueryParams", func(t *testing.T) {
		pathParams := map[string]string{"id": "123", "type": "user"}
		queryParams := map[string]string{"filter": "active", "page": "1"}

		url := client.buildURL("/users/{id}/{type}", pathParams, queryParams)

		expected := fmt.Sprintf("%s://%s/users/123/user?access_token=test-token&filter=active&page=1", config.Scheme, config.Host)
		if url != expected {
			t.Errorf("Expected URL %s, got %s", expected, url)
		}
	})

	t.Run("OnlyPathParams", func(t *testing.T) {
		pathParams := map[string]string{"id": "123"}
		queryParams := map[string]string{}

		url := client.buildURL("/users/{id}", pathParams, queryParams)

		expected := fmt.Sprintf("%s://%s/users/123?access_token=test-token", config.Scheme, config.Host)
		if url != expected {
			t.Errorf("Expected URL %s, got %s", expected, url)
		}
	})

	t.Run("NoParams", func(t *testing.T) {
		pathParams := map[string]string{}
		queryParams := map[string]string{}

		url := client.buildURL("/users", pathParams, queryParams)

		expected := fmt.Sprintf("%s://%s/users?access_token=test-token", config.Scheme, config.Host)
		if url != expected {
			t.Errorf("Expected URL %s, got %s", expected, url)
		}
	})
}

func TestParsePathParams(t *testing.T) {
	client := NewClient("test-token").(*Client)
	pathParams := make(map[string]string)

	field := reflect.StructField{
		Name: "ID",
		Tag:  reflect.StructTag(`path:"userId"`),
	}
	value := reflect.ValueOf("123")

	client.parsePathParams(field, value, pathParams)

	if pathParams["userId"] != "123" {
		t.Errorf("Expected pathParams[userId]=123, got %s", pathParams["userId"])
	}
}

func TestParseQueryParams(t *testing.T) {
	client := NewClient("test-token").(*Client)

	t.Run("RequiredParam", func(t *testing.T) {
		queryParams := make(map[string]string)

		field := reflect.StructField{
			Name: "Filter",
			Tag:  reflect.StructTag(`query:"filter"`),
		}
		value := reflect.ValueOf("active")

		client.parseQueryParams(field, value, queryParams)

		if queryParams["filter"] != "active" {
			t.Errorf("Expected queryParams[filter]=active, got %s", queryParams["filter"])
		}
	})

	t.Run("OmitEmptyParam", func(t *testing.T) {
		queryParams := make(map[string]string)

		field := reflect.StructField{
			Name: "Optional",
			Tag:  reflect.StructTag(`query:"optional,omitempty"`),
		}
		var nilPtr *string
		value := reflect.ValueOf(nilPtr)

		client.parseQueryParams(field, value, queryParams)

		if _, exists := queryParams["optional"]; exists {
			t.Error("Expected optional param to be omitted")
		}
	})
}

func TestParseJSONBody(t *testing.T) {
	client := NewClient("test-token").(*Client)

	t.Run("RequiredField", func(t *testing.T) {
		jsonBody := make(map[string]interface{})

		field := reflect.StructField{
			Name: "Data",
			Tag:  reflect.StructTag(`json:"data"`),
		}
		value := reflect.ValueOf("test_value")

		client.parseJSONBody(field, value, jsonBody)

		if jsonBody["data"] != "test_value" {
			t.Errorf("Expected jsonBody[data]=test_value, got %v", jsonBody["data"])
		}
	})

	t.Run("OmitEmptyField", func(t *testing.T) {
		jsonBody := make(map[string]interface{})

		field := reflect.StructField{
			Name: "Optional",
			Tag:  reflect.StructTag(`json:"optional,omitempty"`),
		}
		var nilPtr *string
		value := reflect.ValueOf(nilPtr)

		client.parseJSONBody(field, value, jsonBody)

		if _, exists := jsonBody["optional"]; exists {
			t.Error("Expected optional field to be omitted")
		}
	})

	t.Run("IgnoredField", func(t *testing.T) {
		jsonBody := make(map[string]interface{})

		field := reflect.StructField{
			Name: "Ignored",
			Tag:  reflect.StructTag(`json:"-"`),
		}
		value := reflect.ValueOf("ignored_value")

		client.parseJSONBody(field, value, jsonBody)

		if _, exists := jsonBody["Ignored"]; exists {
			t.Error("Expected ignored field to not be present")
		}
	})
}
