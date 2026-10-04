package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json/v2"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/agentable/go-intl/tools/gen-cldr/tzdb"
)

func writeTZDataFixture(t *testing.T) string {
	t.Helper()
	var archive bytes.Buffer
	gz := gzip.NewWriter(&archive)
	tw := tar.NewWriter(gz)
	for _, name := range []string{"africa", "antarctica", "asia", "australasia", "europe", "northamerica", "southamerica", "etcetera", "factory", "backward", "backzone", "zone.tab", "version", "LICENSE"} {
		content := "# no definitions\n"
		switch name {
		case "asia":
			content = "Zone Asia/Kolkata 5:30 - IST\nLink Asia/Kolkata Asia/Calcutta\n"
		case "etcetera":
			content = "Zone Etc/UTC 0 - UTC\nLink Etc/UTC UTC\n"
		case "zone.tab":
			content = "IN\t+2232+08822\tAsia/Kolkata\n"
		case "version":
			content = "2025b\n"
		case "LICENSE":
			content = "This data is in the public domain.\n"
		}
		if err := tw.WriteHeader(&tar.Header{Name: name, Mode: 0o600, Size: int64(len(content))}); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, ".tzdata"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".tzdata", "tzdata2025b.tar.gz"), archive.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
	pin := tzdb.Pin{Version: "2025b", URL: "https://data.iana.org/time-zones/releases/tzdata2025b.tar.gz", SHA256: fmt.Sprintf("%x", sha256.Sum256(archive.Bytes())), License: "public-domain"}
	data, err := json.Marshal(pin)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "tzdata.json")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}
