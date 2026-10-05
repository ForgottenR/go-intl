package main

import (
	"context"
	"crypto/sha256"
	"encoding/json/v2"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/agentable/go-intl/tools/internal/datapin"
)

func TestFetchTZData(t *testing.T) {
	t.Parallel()
	const archive = "pinned archive bytes"
	for _, tc := range []struct {
		name, cached, response string
		fail                   bool
	}{
		{"download", "", archive, false}, {"cache hit", archive, "must not download", false},
		{"replace corrupt cache", "old archive", archive, false}, {"bad download keeps cache", "old archive", "bad bytes", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var requests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests.Add(1)
				if _, err := fmt.Fprint(w, tc.response); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			dir := t.TempDir()
			cache := filepath.Join(dir, "cache")
			path := filepath.Join(cache, "tzdata2025b.tar.gz")
			if tc.cached != "" {
				if err := os.MkdirAll(cache, 0o700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte(tc.cached), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			pin := datapin.TZData{Version: "2025b", URL: server.URL + "/archive", SHA256: fmt.Sprintf("%x", sha256.Sum256([]byte(archive))), License: "public-domain"}
			raw, err := json.Marshal(pin)
			if err != nil {
				t.Fatal(err)
			}
			// A valid lock with split fields and escaped slashes must use JSON semantics.
			raw = []byte(strings.ReplaceAll(strings.ReplaceAll(string(raw), `":`, "\":\n"), "/", `\/`))
			lock := filepath.Join(dir, "tzdata.json")
			if err := os.WriteFile(lock, raw, 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := datapin.ReadTZData(lock); err != nil {
				t.Fatal(err)
			}
			err = fetchTZData(context.Background(), lock, cache)
			if (err != nil) != tc.fail {
				t.Fatalf("fetchTZData error = %v, fail=%t", err, tc.fail)
			}
			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			want := archive
			if tc.fail {
				want = tc.cached
			}
			if string(got) != want {
				t.Fatalf("archive = %q, want %q", got, want)
			}
			count := int32(1)
			if tc.name == "cache hit" {
				count = 0
			}
			if requests.Load() != count {
				t.Fatalf("requests=%d,want %d", requests.Load(), count)
			}
			entries, err := os.ReadDir(cache)
			if err != nil || len(entries) != 1 {
				t.Fatalf("cache entries=%v,%v", entries, err)
			}
		})
	}
}
