package status_icons

import (
	"bytes"
	"image"
	_ "image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestEmbeddedIconsAreValidPNG(t *testing.T) {
	icons := []struct {
		name string
		data []byte
	}{
		{"default", defaultIconData},
		{"crash", crashIconData},
		{"pause", pauseIconData},
		{"live-reload", liveReloadIconData},
	}
	for _, ic := range icons {
		img, format, err := image.Decode(bytes.NewReader(ic.data))
		if err != nil {
			t.Errorf("image.Decode(%s) error = %v", ic.name, err)
			continue
		}
		if format != "png" {
			t.Errorf("image format for %s = %s, want png", ic.name, format)
		}
		if img.Bounds().Dx() == 0 || img.Bounds().Dy() == 0 {
			t.Errorf("image %s has empty dimensions", ic.name)
		}
	}
}

func resetDefaultIcons(t *testing.T) {
	origDefault := Default
	origCrash := Crash
	origPause := Pause
	origLiveReload := LiveReload
	t.Cleanup(func() {
		Default = origDefault
		Crash = origCrash
		Pause = origPause
		LiveReload = origLiveReload
	})
}

func TestLoadCustomStatusIcons(t *testing.T) {
	resetDefaultIcons(t)

	configDir := t.TempDir()
	iconsDir := filepath.Join(configDir, statusIconsDir)
	if err := os.MkdirAll(iconsDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Use valid PNG data (reusing liveReloadIconData and crashIconData) as custom icons.
	files := map[string][]byte{
		"default.png":      liveReloadIconData,
		"pause.custom.png": crashIconData,
		"crash.png":        {},
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(iconsDir, name), content, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	crashBefore := Crash
	if err := LoadCustomStatusIcons(configDir); err != nil {
		t.Fatalf("LoadCustomStatusIcons: %v", err)
	}

	if bytes.Equal(Default.Data, defaultIconData) {
		t.Error("Default was not updated by custom default.png")
	}
	if !bytes.Equal(Default.Data, liveReloadIconData) {
		t.Error("Default.Data does not match custom default.png content")
	}

	if bytes.Equal(Pause.Data, pauseIconData) {
		t.Error("Pause was not updated by custom pause.custom.png")
	}
	if !bytes.Equal(Pause.Data, crashIconData) {
		t.Error("Pause.Data does not match custom pause.custom.png content")
	}

	// An empty file is ignored instead of being passed to systray.
	if !bytes.Equal(Crash.Data, crashBefore.Data) {
		t.Error("Crash was overwritten by an empty icon file")
	}

	// Not overridden, keeps the embedded icon.
	if len(LiveReload.Data) == 0 {
		t.Errorf("LiveReload = %d bytes, want embedded icon", len(LiveReload.Data))
	}
}

func TestLoadCustomStatusIconsIgnoresNonPNG(t *testing.T) {
	resetDefaultIcons(t)

	configDir := t.TempDir()
	iconsDir := filepath.Join(configDir, statusIconsDir)
	if err := os.MkdirAll(iconsDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Legacy ICO files left by previous kanata-tray versions
	if err := os.WriteFile(filepath.Join(iconsDir, "default.ico"), []byte("legacy ico bytes"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Invalid PNG file (wrong format)
	if err := os.WriteFile(filepath.Join(iconsDir, "crash.png"), []byte("not a real png"), 0o644); err != nil {
		t.Fatal(err)
	}

	defaultBefore := Default
	crashBefore := Crash
	if err := LoadCustomStatusIcons(configDir); err != nil {
		t.Fatalf("LoadCustomStatusIcons: %v", err)
	}

	if !bytes.Equal(Default.Data, defaultBefore.Data) {
		t.Errorf("Default was incorrectly overwritten by non-PNG file (default.ico)")
	}
	if !bytes.Equal(Crash.Data, crashBefore.Data) {
		t.Errorf("Crash was incorrectly overwritten by invalid PNG file (crash.png)")
	}
}

func TestLoadCustomStatusIconsMissingDir(t *testing.T) {
	if err := LoadCustomStatusIcons(t.TempDir()); err != nil {
		t.Errorf("LoadCustomStatusIcons on a config dir without status_icons: %v", err)
	}
}

func TestCreateDefaultStatusIconsDirIfNotExists(t *testing.T) {
	configDir := t.TempDir()
	iconsDir := filepath.Join(configDir, statusIconsDir)

	// Pre-create the directory with a legacy ICO and an existing custom default.png
	if err := os.MkdirAll(iconsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	customPNG := []byte("custom-png-marker")
	if err := os.WriteFile(filepath.Join(iconsDir, "default.png"), customPNG, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(iconsDir, "legacy.ico"), []byte("ico"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := CreateDefaultStatusIconsDirIfNotExists(configDir); err != nil {
		t.Fatalf("CreateDefaultStatusIconsDirIfNotExists: %v", err)
	}

	// Existing default.png should not be overwritten
	gotDefault, err := os.ReadFile(filepath.Join(iconsDir, "default.png"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(gotDefault, customPNG) {
		t.Errorf("existing default.png was overwritten")
	}

	// Missing PNGs should have been created
	for _, name := range []string{"crash.png", "pause.png", "live-reload.png"} {
		path := filepath.Join(iconsDir, name)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("missing icon %s was not created: %v", name, err)
		}
	}
}
