package datapin

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json/v2"
	"fmt"
	"net/url"
	"os"
	"strings"
)

// TZData is the pinned archive identity; it does not fetch or inspect the archive.
type TZData struct {
	Version string `json:"version"`
	URL     string `json:"url"`
	SHA256  string `json:"sha256"`
	License string `json:"license"`
}

// ReadTZData validates the shared lock syntax and canonicalizes the hash to lowercase.
func ReadTZData(path string) (TZData, error) {
	//nolint:gosec // G304: maintainer tools explicitly select the pin file.
	data, err := os.ReadFile(path)
	if err != nil {
		return TZData{}, fmt.Errorf("read %s: %w", path, err)
	}
	var pin TZData
	if err := json.Unmarshal(data, &pin); err != nil {
		return TZData{}, fmt.Errorf("parse %s: %w", path, err)
	}
	if len(pin.Version) < 5 || !asciiDigits(pin.Version[:4]) || strings.IndexFunc(pin.Version[4:], func(r rune) bool { return r < 'a' || r > 'z' }) >= 0 {
		return TZData{}, fmt.Errorf("%s: invalid or missing tzdata version %q", path, pin.Version)
	}
	hash, err := hex.DecodeString(pin.SHA256)
	if err != nil || len(hash) != sha256.Size {
		return TZData{}, fmt.Errorf("%s: invalid tzdata sha256 %q", path, pin.SHA256)
	}
	u, err := url.Parse(pin.URL)
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") || u.Hostname() == "" || u.User != nil || u.Fragment != "" {
		return TZData{}, fmt.Errorf("%s: invalid tzdata URL %q", path, pin.URL)
	}
	if pin.License != "public-domain" {
		return TZData{}, fmt.Errorf("%s: unsupported tzdata license %q", path, pin.License)
	}
	pin.SHA256 = strings.ToLower(pin.SHA256)
	return pin, nil
}

func asciiDigits(text string) bool {
	if text == "" {
		return false
	}
	for _, r := range text {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
