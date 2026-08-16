package articles

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testHandler(t *testing.T) http.Handler {
	t.Helper()
	data, err := LoadYAML(Fixture)
	if err != nil {
		t.Fatal(err)
	}
	return NewHandler(NewService(data), http.NotFoundHandler())
}

func TestSearchEndpointReturnsJSONArrayForNoMatch(t *testing.T) {
	server := httptest.NewServer(testHandler(t))
	defer server.Close()
	response, err := http.Get(server.URL + "/api/articles?q=没有匹配的布线关键词")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.StatusCode)
	}
	var articles []Article
	if err := json.Unmarshal(body, &articles); err != nil {
		t.Fatalf("invalid article list: %v", err)
	}
	if articles == nil {
		t.Fatalf("expected [] in response, got %s", strings.TrimSpace(string(body)))
	}
	if len(articles) != 0 {
		t.Fatalf("expected no results, got %d", len(articles))
	}
}

func TestArticleDetailContainsLearningSections(t *testing.T) {
	server := httptest.NewServer(testHandler(t))
	defer server.Close()
	response, err := http.Get(server.URL + "/api/articles/cat6-cable")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var article Article
	if err := json.NewDecoder(response.Body).Decode(&article); err != nil {
		t.Fatal(err)
	}
	if article.Title == "" || len(article.Materials) == 0 || len(article.Steps) == 0 || len(article.Troubleshooting) == 0 {
		t.Fatalf("incomplete article: %#v", article)
	}
}

func TestToolsEndpointIsAvailable(t *testing.T) {
	server := httptest.NewServer(testHandler(t))
	defer server.Close()
	response, err := http.Get(server.URL + "/api/tools")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var tools []map[string]string
	if err := json.NewDecoder(response.Body).Decode(&tools); err != nil {
		t.Fatal(err)
	}
	if len(tools) < 3 {
		t.Fatalf("expected tool list, got %d", len(tools))
	}
}
