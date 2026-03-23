package jsonfields_test

import (
	"encoding/json"
	"testing"

	"github.com/wangke19/po/pkg/jsonfields"
)

func TestFilterFields_subset(t *testing.T) {
	type item struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Extra string `json:"extra"`
	}
	v := item{ID: "WI-1", Title: "My test", Extra: "ignored"}
	got, err := jsonfields.FilterFields(v, []string{"id", "title"})
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(got, &m); err != nil {
		t.Fatal(err)
	}
	if _, ok := m["extra"]; ok {
		t.Error("extra field should be filtered out")
	}
	if m["id"] != "WI-1" {
		t.Errorf("id: got %v", m["id"])
	}
}

func TestFilterFields_empty_returns_all(t *testing.T) {
	v := map[string]any{"a": 1, "b": 2}
	got, err := jsonfields.FilterFields(v, nil)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(got, &m); err != nil {
		t.Fatal(err)
	}
	if len(m) != 2 {
		t.Errorf("expected 2 keys, got %d", len(m))
	}
}

func TestFilterFields_slice(t *testing.T) {
	type item struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}
	items := []item{{ID: "1", Title: "a"}, {ID: "2", Title: "b"}}
	got, err := jsonfields.FilterFields(items, []string{"id"})
	if err != nil {
		t.Fatal(err)
	}
	var result []map[string]any
	if err := json.Unmarshal(got, &result); err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 items, got %d", len(result))
	}
	if _, ok := result[0]["title"]; ok {
		t.Error("title should be filtered out")
	}
}
