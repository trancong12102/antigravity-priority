package host

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// ErrDocumentChanged reports that the credential document on disk no longer
// matches the version a replacement was computed from, so the replacement was
// not written.
var ErrDocumentChanged = errors.New("auth document changed since it was read")

// DocumentPatcher is the low-level complete-document replacement primitive.
// Field-specific patch methods intentionally do not exist: the Apply-layer
// Host Transition owns expected-state checks and constructs the full document
// before calling this primitive.
type DocumentPatcher struct{}

// NewDocumentPatcher creates a complete-document replacement primitive.
func NewDocumentPatcher() *DocumentPatcher {
	return &DocumentPatcher{}
}

// Replace atomically replaces the physical credential document at path.
func (p *DocumentPatcher) Replace(ctx context.Context, path string, data []byte) error {
	trimmedPath := strings.TrimSpace(path)
	if trimmedPath == "" {
		return fmt.Errorf("%w: auth document path is required", ErrInvalidRequest)
	}
	return writeFileAtomic(ctx, trimmedPath, data, nil)
}

// writeFileAtomic replaces path with data through a temp file and rename. When
// base is non-nil the file is re-read right before the rename and the write is
// abandoned with ErrDocumentChanged unless it still holds the same JSON
// document as base.
func writeFileAtomic(ctx context.Context, path string, data, base []byte) (err error) {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("write auth document context: %w", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat auth document: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("write auth document context: %w", err)
	}
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, filepath.Base(path)+".tmp-*")
	if err != nil {
		return fmt.Errorf("create auth document temp: %w", err)
	}
	tmpName := tmp.Name()
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.Remove(tmpName)
		}
	}()
	if err := ctx.Err(); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write auth document context: %w", err)
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write auth document temp: %w", err)
	}
	if err := ctx.Err(); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write auth document context: %w", err)
	}
	if err := tmp.Chmod(info.Mode().Perm()); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("chmod auth document temp: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close auth document temp: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("write auth document context: %w", err)
	}
	if base != nil {
		current, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("re-read auth document: %w", err)
		}
		if !sameJSONDocument(current, base) {
			return ErrDocumentChanged
		}
	}
	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("replace auth document: %w", err)
	}
	cleanup = false
	return nil
}

// sameJSONDocument compares two JSON documents by value. host.auth.get returns
// the file re-encoded by CPA's RPC layer, so a byte comparison against the file
// would report a change on every pretty-printed or HTML-escaped document.
// Anything that does not parse, such as a file caught mid-write, is different.
func sameJSONDocument(a, b []byte) bool {
	if bytes.Equal(a, b) {
		return true
	}
	left, ok := decodeJSONValue(a)
	if !ok {
		return false
	}
	right, ok := decodeJSONValue(b)
	if !ok {
		return false
	}
	return reflect.DeepEqual(left, right)
}

func decodeJSONValue(raw []byte) (any, bool) {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, false
	}
	if decoder.More() {
		return nil, false
	}
	return value, true
}
