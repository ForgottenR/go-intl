package testcontract

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"reflect"
	"testing"
)

// AssertResolvedOptionsJSON compares the complete marshaled record, including
// property presence and JSON types. Object property order has no significance.
func AssertResolvedOptionsJSON(t testing.TB, got any, want jsontext.Value) {
	t.Helper()

	data, err := json.Marshal(got, json.Deterministic(true))
	if err != nil {
		t.Fatalf("ResolvedOptions() JSON: %v", err)
	}
	var actual, expected map[string]any
	if jsontext.Value(data).Kind() != '{' || want.Kind() != '{' {
		t.Fatalf("ResolvedOptions() must compare objects: got %s, want %s", data, want)
	}
	if err := json.Unmarshal(data, &actual); err != nil {
		t.Fatalf("ResolvedOptions() JSON: %v", err)
	}
	if err := json.Unmarshal(want, &expected); err != nil {
		t.Fatalf("expectedResolvedOptions JSON: %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("ResolvedOptions() JSON = %s, want %s", data, want)
	}
}
