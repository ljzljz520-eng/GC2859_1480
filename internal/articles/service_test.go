package articles

import (
	"sync"
	"testing"
)

func TestSearchReturnsKnownArticle(t *testing.T) {
	data, err := LoadYAML(Fixture)
	if err != nil {
		t.Fatal(err)
	}
	got := NewService(data).Search("水晶头")
	if len(got) != 1 || got[0].ID != "rj45-crimp" {
		t.Fatalf("unexpected search result: %#v", got)
	}
}

func TestSearchNoMatchReturnsEmptyArrayValue(t *testing.T) {
	data, err := LoadYAML(Fixture)
	if err != nil {
		t.Fatal(err)
	}
	got := NewService(data).Search("不存在的布线关键词")
	if got == nil {
		t.Fatal("expected a non-nil empty result")
	}
	if len(got) != 0 {
		t.Fatalf("expected no results, got %d", len(got))
	}
}

func TestSearchCanBeCalledConcurrently(t *testing.T) {
	data, err := LoadYAML(Fixture)
	if err != nil {
		t.Fatal(err)
	}
	service := NewService(data)
	start := make(chan struct{})
	results := make(chan int, 3)
	var ready sync.WaitGroup
	ready.Add(3)
	for _, query := range []string{"水晶头", "T568B", "PoE"} {
		go func(query string) {
			ready.Done()
			<-start
			results <- len(service.Search(query))
		}(query)
	}
	ready.Wait()
	close(start)
	for range 3 {
		if n := <-results; n != 1 {
			t.Fatalf("expected one result, got %d", n)
		}
	}
}
