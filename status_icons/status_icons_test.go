package status_icons

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCustomStatusIcons(t *testing.T) {
	configDir := t.TempDir()
	iconsDir := filepath.Join(configDir, statusIconsDir)
	if err := os.MkdirAll(iconsDir, 0o755); err != nil {
		t.Fatal(err)
	}

	files := map[string]string{
		"default.png":      "custom default",
		"pause.custom.png": "custom pause",
		"crash.png":        "",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(iconsDir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	crashBefore := Crash
	if err := LoadCustomStatusIcons(configDir); err != nil {
		t.Fatalf("LoadCustomStatusIcons: %v", err)
	}

	if got := string(Default.Data); got != "custom default" {
		t.Errorf("Default.Data = %q, want custom default", got)
	}

	if got := string(Pause.Data); got != "custom pause" {
		t.Errorf("Pause.Data = %q, want custom pause", got)
	}

	// An empty file is ignored instead of being passed to systray.
	if string(Crash.Data) != string(crashBefore.Data) {
		t.Error("Crash was overwritten by an empty icon file")
	}

	// Not overridden, keeps the embedded icon.
	if len(LiveReload.Data) == 0 {
		t.Errorf("LiveReload = %d bytes, want embedded icon", len(LiveReload.Data))
	}
}

func TestLoadCustomStatusIconsMissingDir(t *testing.T) {
	if err := LoadCustomStatusIcons(t.TempDir()); err != nil {
		t.Errorf("LoadCustomStatusIcons on a config dir without status_icons: %v", err)
	}
}
