package golden

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google"
)

type Harness struct {
	*testing.T
	Config       *google.Config
	Resources    map[string]*schema.Resource
	RoundTripper *mockRoundTripper
}

func NewHarness(t *testing.T) *Harness {
	rm, err := google.ResourceMapWithErrors()
	if err != nil {
		t.Fatalf("ResourceMapWithErrors failed: %v", err)
	}

	config := &google.Config{
		AccessToken:         "coolbeans",
		Project:             "config-test",
		BillingProject:      "billing-project",
		UserProjectOverride: true,
		Region:              "us-central1",
	}

	roundTripper := &mockRoundTripper{}

	config.HTTPClient = &http.Client{
		Transport: roundTripper,
	}

	google.ConfigureBasePaths(config)

	if err := config.LoadAndValidate(context.Background()); err != nil {
		t.Fatalf("config.LoadAndValidate failed: %v", err)
	}
	return &Harness{
		T:            t,
		Resources:    rm,
		Config:       config,
		RoundTripper: roundTripper,
	}
}

func (h *Harness) MustHaveXGoogUserProject(project string) {
	for _, request := range h.RoundTripper.Requests {
		h.Logf("request: %#v", request)
		if got, want := request.Header.Get("X-Goog-User-Project"), project; got != want {
			h.Errorf("unexpected X-Goog-User-Project; got %q, want %q", got, want)
		}
	}
}
func TestResourceSQLDatabase(t *testing.T) {
	h := NewHarness(t)

	resourceID := "google_sql_database_instance"
	resource := h.Resources[resourceID]
	if resource == nil {
		t.Fatalf("resource %q not found", resourceID)
	}

	{
		data := resource.TestResourceData()
		data.Set("name", "test1")

		if err := resource.Create(data, h.Config); err != nil {
			t.Errorf("Create failed: %v", err)
		}

		h.MustHaveXGoogUserProject("billing-project")
	}

	{
		data := resource.TestResourceData()
		data.Set("name", "test1")

		if err := resource.Read(data, h.Config); err != nil {
			t.Errorf("Read failed: %v", err)
		}

		h.MustHaveXGoogUserProject("billing-project")
	}
}

func TestResourceContainerCluster(t *testing.T) {
	h := NewHarness(t)

	resourceID := "google_container_cluster"
	resource := h.Resources[resourceID]
	if resource == nil {
		t.Fatalf("resource %q not found", resourceID)
	}

	{
		data := resource.TestResourceData()
		data.Set("name", "test1")
		data.Set("project", "resource-project")
		data.Set("location", "us-central1-a")

		if err := resource.Create(data, h.Config); err != nil {
			t.Errorf("Create failed: %v", err)
		}

		h.MustHaveXGoogUserProject("resource-project")
	}

	{
		data := resource.TestResourceData()
		data.Set("name", "test1")
		data.Set("project", "resource-project")
		data.Set("location", "us-central1-a")

		if err := resource.Read(data, h.Config); err != nil {
			t.Errorf("Read failed: %v", err)
		}

		h.MustHaveXGoogUserProject("resource-project")
	}
}

type Request struct {
	Method string
	URL    string
	Header http.Header
	Body   string
}
type mockRoundTripper struct {
	Requests []Request
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	//log.Printf("request: %#v", req)
	log.Printf("request: %v %v", req.Method, req.URL)

	request := fmt.Sprintf("%s %s", req.Method, req.URL)
	body := make(map[string]interface{})

	{
		c := Request{
			Method: req.Method,
			URL:    req.URL.String(),
			Header: req.Header,
		}

		if req.Body != nil {
			requestBody, err := ioutil.ReadAll(req.Body)
			if err != nil {
				panic("failed to read request body")
			}
			c.Body = string(requestBody)
		}
		m.Requests = append(m.Requests, c)
	}

	response := &http.Response{
		StatusCode: 403,
		Status:     "mockRoundTripper injecting fake response",
	}

	if request == "GET https://openidconnect.googleapis.com/v1/userinfo?alt=json" {
		body["email"] = "test@example.com"

		response.StatusCode = 200
	}

	if request == "POST https://sqladmin.googleapis.com/sql/v1beta4/projects/config-test/instances?alt=json&prettyPrint=false" {
		response.StatusCode = 200
	}

	if request == "GET https://sqladmin.googleapis.com/sql/v1beta4/projects/config-test/operations/?alt=json&prettyPrint=false" {
		body["status"] = "DONE"

		response.StatusCode = 200
	}

	if request == "GET https://sqladmin.googleapis.com/sql/v1beta4/projects/config-test/instances/test1?alt=json&prettyPrint=false" {
		body["settings"] = map[string]interface{}{}

		response.StatusCode = 200
	}

	if request == "GET https://sqladmin.googleapis.com/sql/v1beta4/projects/config-test/instances/test1/users?alt=json&prettyPrint=false" {
		response.StatusCode = 200
	}

	if request == "GET https://container.googleapis.com/v1beta1/projects/resource-project/locations/us-central1-a/clusters/test1?alt=json&prettyPrint=false" {
		body["legacyAbac"] = map[string]interface{}{}
		body["networkConfig"] = map[string]interface{}{}
		body["status"] = "RUNNING"
		response.StatusCode = 200
	}

	if request == "POST https://container.googleapis.com/v1beta1/projects/resource-project/locations/us-central1-a/clusters?alt=json&prettyPrint=false" {
		response.StatusCode = 200
	}

	if request == "GET https://container.googleapis.com/v1beta1/projects/resource-project/locations/us-central1-a/operations/?alt=json&prettyPrint=false" {
		body["status"] = "DONE"

		response.StatusCode = 200
	}

	if body != nil {
		j, err := json.Marshal(body)
		if err != nil {
			panic("json.Marshal failed")
		}

		log.Printf("response: %d %s", response.StatusCode, string(j))

		response.Body = ioutil.NopCloser(bytes.NewReader(j))
	} else {
		log.Printf("response: %d %s", response.StatusCode, "-")
	}

	return response, nil
}
