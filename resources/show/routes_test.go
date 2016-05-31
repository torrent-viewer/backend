package show

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	// Initialize SQLite driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/shwoodard/jsonapi"
	"github.com/torrent-viewer/backend/datastore"
	"github.com/torrent-viewer/backend/router"
)

var (
	server          *httptest.Server
	baseURL         string
	integerOverflow string = "9223372036854775808"
)

func TestMain(m *testing.M) {
	flag.Parse()
	datastore.Init("sqlite3", "", "", "", "", "/tmp/torrent-viewer-test.db")
	datastore.Conn.AutoMigrate(&Show{})
	r := router.NewRouter()
	r.AddResource("shows", ShowResource{})
	server = httptest.NewServer(r)
	baseURL = fmt.Sprintf("%s/shows", server.URL)
	ret := m.Run()
	datastore.Conn.DropTable(&Show{})
	os.Exit(ret)
}

func TestShowsIndex(t *testing.T) {
	request, err := http.NewRequest("GET", baseURL, nil)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	if response.StatusCode != 200 {
		t.Errorf("Expected HTTP 200, got HTTP %d", response.StatusCode)
	}
}

func testEndpoint(t *testing.T, method string, url string, input *string) *http.Response {
	var reader io.Reader
	if input != nil {
		reader = strings.NewReader(*input)
	} else {
		reader = nil
	}
	request, err := http.NewRequest(method, url, reader)
	if err != nil {
		t.Error(err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	return response
}

func TestShowsStore(t *testing.T) {
	var input string

	input = `{"id": 5}`
	response := testEndpoint(t, "POST", baseURL, &input)
	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP %d, got HTTP %d", http.StatusBadRequest, response.StatusCode)
	}
	input = `{
    "data": {
      "type": "shows",
      "attributes": {
        "year": 1996
      }
    }
  }`
	response = testEndpoint(t, "POST", baseURL, &input)
	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP %d, got HTTP %d", http.StatusBadRequest, response.StatusCode)
	}
	input = `{
    "data": {
      "type": "shows",
      "attributes": {
        "title": "Star Wars VII",
        "year": 2015
      }
    }
  }`
	response = testEndpoint(t, "POST", baseURL, &input)
	if response.StatusCode != http.StatusCreated {
		t.Errorf("Expected HTTP %d, got HTTP %d", http.StatusCreated, response.StatusCode)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	input = buf.String()
	response = testEndpoint(t, "POST", baseURL, &input)
	if response.StatusCode != http.StatusConflict {
		t.Errorf("Expected HTTP %d, got HTTP %d", http.StatusConflict, response.StatusCode)
	}
}

func TestShowsShow(t *testing.T) {
	input := `{
    "data": {
      "type": "shows",
      "attributes": {
        "title": "Star Wars VII",
        "year": 2015
      }
    }
  }`
	response := testEndpoint(t, "POST", baseURL, &input)
	if response.StatusCode != http.StatusCreated {
		t.Errorf("Expected HTTP %d, got HTTP %d", http.StatusCreated, response.StatusCode)
	}
	response = testEndpoint(t, "GET", fmt.Sprintf("%s/%d", baseURL, 1), nil)
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP %d, got HTTP %d", http.StatusOK, response.StatusCode)
	}
	response = testEndpoint(t, "GET", fmt.Sprintf("%s/%d", baseURL, math.MaxInt32), nil)
	if response.StatusCode != http.StatusNotFound {
		t.Errorf("Expected HTTP %d, got HTTP %d", http.StatusNotFound, response.StatusCode)
	}
	response = testEndpoint(t, "GET", fmt.Sprintf("%s/%d", baseURL, math.MaxInt32), nil)
	if response.StatusCode != http.StatusNotFound {
		t.Errorf("Expected HTTP %d, got HTTP %d", http.StatusNotFound, response.StatusCode)
	}
	response = testEndpoint(t, "GET", fmt.Sprintf("%s/%s", baseURL, integerOverflow), nil)
	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP %d, got HTTP %d", http.StatusBadRequest, response.StatusCode)
	}
}

func TestShowsUpdate(t *testing.T) {
	response := testEndpoint(t, "GET", fmt.Sprintf("%s/%s", baseURL, integerOverflow), nil)
	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP %d, got HTTP %d", http.StatusBadRequest, response.StatusCode)
	}
	input := `{
    "data": {
      "type": "shows",
      "attributes": {
        "title": "Star Wars VII",
        "year": 2015
      }
    }
  }`
	response = testEndpoint(t, "POST", baseURL, &input)
	if response.StatusCode != http.StatusCreated {
		t.Errorf("Expected HTTP %d, got HTTP %d", http.StatusCreated, response.StatusCode)
	}
	var show Show
	if err := jsonapi.UnmarshalPayload(response.Body, &show); err != nil {
		t.Error(err)
		return
	}
	show.Title = "GhostBusters"
	buf := new(bytes.Buffer)
	if err := jsonapi.MarshalOnePayload(buf, &show); err != nil {
		t.Error(err)
		return
	}
	input = buf.String()
	response = testEndpoint(t, "PATCH", fmt.Sprintf("%s/%d", baseURL, show.ID+1), &input)
	if response.StatusCode != http.StatusNotFound {
		t.Errorf("Expected HTTP %d, got HTTP %d", http.StatusNotFound, response.StatusCode)
	}
	response = testEndpoint(t, "PATCH", fmt.Sprintf("%s/%s", baseURL, integerOverflow), &input)
	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP %d, got HTTP %d", http.StatusBadRequest, response.StatusCode)
	}
	response = testEndpoint(t, "PATCH", fmt.Sprintf("%s/%d", baseURL, show.ID), &input)
	if response.StatusCode != http.StatusNoContent {
		t.Errorf("Expected HTTP %d, got HTTP %d", http.StatusCreated, response.StatusCode)
	}
	show.ID += 1000
	buf = new(bytes.Buffer)
	if err := jsonapi.MarshalOnePayload(buf, &show); err != nil {
		t.Error(err)
		return
	}
	input = buf.String()
	response = testEndpoint(t, "PATCH", fmt.Sprintf("%s/%d", baseURL, show.ID - 1000), &input)
	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP %d, got HTTP %d", http.StatusBadRequest, response.StatusCode)
	}
}

func TestShowsDestroy(t *testing.T) {
	input := `{
    "data": {
      "type": "shows",
      "attributes": {
        "title": "Star Wars VII",
        "year": 2015
      }
    }
  }`
	response := testEndpoint(t, "POST", baseURL, &input)
	if response.StatusCode != http.StatusCreated {
		t.Errorf("Expected HTTP %d, got HTTP %d", http.StatusCreated, response.StatusCode)
	}
	var show Show
	if err := jsonapi.UnmarshalPayload(response.Body, &show); err != nil {
		t.Error(err)
		return
	}
	response = testEndpoint(t, "DELETE", fmt.Sprintf("%s/%s", baseURL, integerOverflow), nil)
	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP %d, got HTTP %d", http.StatusBadRequest, response.StatusCode)
	}
	response = testEndpoint(t, "DELETE", fmt.Sprintf("%s/%d", baseURL, math.MaxInt32), nil)
	if response.StatusCode != http.StatusNotFound {
		t.Errorf("Expected HTTP %d, got HTTP %d", http.StatusNotFound, response.StatusCode)
	}
	response = testEndpoint(t, "DELETE", fmt.Sprintf("%s/%d", baseURL, show.ID), nil)
	if response.StatusCode != http.StatusNoContent {
		t.Errorf("Expected HTTP %d, got HTTP %d", http.StatusNoContent, response.StatusCode)
	}
}
