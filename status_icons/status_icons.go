package status_icons

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/gommon/log"
)

// Icon is a tray icon.
type Icon struct {
	Data []byte
}

//go:embed default.png
var defaultIconData []byte

//go:embed crash.png
var crashIconData []byte

//go:embed pause.png
var pauseIconData []byte

//go:embed live-reload.png
var liveReloadIconData []byte

var (
	Default    = Icon{Data: defaultIconData}
	Crash      = Icon{Data: crashIconData}
	Pause      = Icon{Data: pauseIconData}
	LiveReload = Icon{Data: liveReloadIconData}
)

//////////////////////////////////////////////

var statusIconsDir string = "status_icons"

func filenameWithoutExt(path string) string {
	basename := filepath.Base(path)
	return strings.TrimSuffix(basename, filepath.Ext(basename))
}

func LoadCustomStatusIcons(configDir string) error {
	dir := filepath.Join(configDir, statusIconsDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("os.ReadDir: %v", err)
	}

	targets := []struct {
		prefix string
		icon   *Icon
	}{
		{"default", &Default},
		{"crash", &Crash},
		{"pause", &Pause},
		{"live-reload", &LiveReload},
	}

	for _, target := range targets {
		filename, ok := findStatusIcon(entries, target.prefix)
		if !ok {
			continue
		}
		path := filepath.Join(dir, filename)

		log.Infof("loading status icon: %s", path)
		fileContent, err := os.ReadFile(path)
		if err != nil {
			log.Errorf("LoadCustomStatusIcons: os.ReadFile: %v", err)
			continue
		}
		if len(fileContent) == 0 {
			log.Errorf("LoadCustomStatusIcons: status icon file is empty: %s", path)
			continue
		}
		if _, err := png.DecodeConfig(bytes.NewReader(fileContent)); err != nil {
			log.Warnf("LoadCustomStatusIcons: %s is not a valid PNG image (%v), keeping default icon", path, err)
			continue
		}

		*target.icon = Icon{Data: fileContent}
	}

	return nil
}

// Finds a status icon file matching a prefix.
// Only PNG files are matched to prevent loading legacy .ico files.
// Only first match is returned, others are ignored.
func findStatusIcon(entries []os.DirEntry, prefix string) (filename string, found bool) {
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.ToLower(filepath.Ext(entry.Name())) != ".png" {
			continue
		}
		name := filenameWithoutExt(entry.Name())
		if strings.HasPrefix(name, prefix) {
			return entry.Name(), true
		}
	}
	return "", false
}

func CreateDefaultStatusIconsDirIfNotExists(configDir string) error {
	customIconsPath := filepath.Join(configDir, statusIconsDir)
	err := os.MkdirAll(customIconsPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create folder: %v", err)
	}

	names := []string{"default.png", "crash.png", "pause.png", "live-reload.png"}
	data := [][]byte{defaultIconData, crashIconData, pauseIconData, liveReloadIconData}
	for i, name := range names {
		path := filepath.Join(customIconsPath, name)
		if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
			if err := os.WriteFile(path, data[i], 0o644); err != nil {
				return fmt.Errorf("writing file %s failed", path)
			}
		} else if err != nil {
			return fmt.Errorf("error checking if %s exists: %v", path, err)
		}
	}
	return nil
}
