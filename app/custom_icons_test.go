package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rszyma/kanata-tray/config"
)

func TestResolveIconsTemplateSuffix(t *testing.T) {
	configDir := t.TempDir()
	iconsDir := filepath.Join(configDir, "icons")
	if err := os.MkdirAll(iconsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"mouseTemplate.png": "template icon",
		"qwerty.png":        "plain icon",
		"allTemplate.png":   "template wildcard",
		"empty.png":         "",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(iconsDir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	cfg := &config.Config{
		PresetDefaults: config.Preset{
			LayerIcons: map[string]string{
				"qwerty": "qwerty.png",
				"*":      "allTemplate.png",
			},
		},
		Presets: config.NewOrderedMap[string, *config.Preset](),
	}
	cfg.Presets.Set("main", &config.Preset{
		LayerIcons: map[string]string{
			"mouse":  "mouseTemplate.png",
			"broken": "empty.png",
		},
	})

	icons := ResolveIcons(configDir, cfg)

	mouse := icons.IconForLayerName("main", "mouse")
	if mouse == nil || string(mouse.Data) != "template icon" {
		t.Errorf("mouse icon = %+v, want the template icon", mouse)
	}

	qwerty := icons.IconForLayerName("main", "qwerty")
	if qwerty == nil || string(qwerty.Data) != "plain icon" {
		t.Errorf("qwerty icon = %+v, want plain icon", qwerty)
	}

	// Falls through to the wildcard.
	other := icons.IconForLayerName("main", "some-other-layer")
	if other == nil || string(other.Data) != "template wildcard" {
		t.Errorf("wildcard icon = %+v, want the template wildcard", other)
	}

	// An empty icon file is skipped at load time, so the layer falls back to
	// the wildcard instead of resolving to unusable icon data.
	broken := icons.IconForLayerName("main", "broken")
	if broken == nil || string(broken.Data) != "template wildcard" {
		t.Errorf("broken icon = %+v, want fallback to the wildcard", broken)
	}

	// No preset entry and no matching layer: unknown preset still gets the
	// global wildcard.
	unknown := icons.IconForLayerName("no-such-preset", "mouse")
	if unknown == nil || string(unknown.Data) != "template wildcard" {
		t.Errorf("unknown preset icon = %+v, want the global wildcard", unknown)
	}
}
