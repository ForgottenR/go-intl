// Command fetch-tzdata fetches and verifies the pinned IANA source archive.
package main

import (
	"context"
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/agentable/go-intl/tools/internal/datapin"
)

func main() {
	lock := flag.String("lock", "tools/gen-cldr/tzdata.json", "archive lock JSON")
	cache := flag.String("cache", "tools/gen-cldr/.tzdata", "archive cache directory")
	flag.Parse()
	if err := fetchTZData(context.Background(), *lock, *cache); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func fetchTZData(ctx context.Context, lock, cache string) error {
	pin, err := datapin.ReadTZData(lock)
	if err != nil {
		return err
	}
	path := filepath.Join(cache, "tzdata"+pin.Version+".tar.gz")
	//nolint:gosec // G304: the maintainer selects the cache directory; the pin version is validated.
	data, err := os.ReadFile(path)
	if err == nil && fmt.Sprintf("%x", sha256.Sum256(data)) == pin.SHA256 {
		return nil
	}
	if err := os.MkdirAll(cache, 0o750); err != nil {
		return fmt.Errorf("create archive cache %s: %w", cache, err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pin.URL, nil)
	if err != nil {
		return fmt.Errorf("request %s: %w", pin.URL, err)
	}
	client := &http.Client{Timeout: time.Minute}
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("fetch %s: %w", pin.URL, err)
	}
	// Closing a read-only response does not affect the verified archive bytes.
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("fetch %s: HTTP %s", pin.URL, response.Status)
	}
	temp, err := os.CreateTemp(cache, ".tzdata-*")
	if err != nil {
		return fmt.Errorf("create archive temporary file in %s: %w", cache, err)
	}
	defer func() { _ = os.Remove(temp.Name()) }()
	hash := sha256.New()
	_, copyErr := io.Copy(io.MultiWriter(temp, hash), response.Body)
	closeErr := temp.Close()
	if copyErr != nil {
		return fmt.Errorf("copy %s to %s: %w", pin.URL, temp.Name(), copyErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close %s: %w", temp.Name(), closeErr)
	}
	if got := fmt.Sprintf("%x", hash.Sum(nil)); got != pin.SHA256 {
		return fmt.Errorf("fetch %s: archive sha256 %s, want %s from %s", pin.URL, got, pin.SHA256, lock)
	}
	if err := os.Rename(temp.Name(), path); err != nil {
		return fmt.Errorf("replace archive %s: %w", path, err)
	}
	return nil
}
