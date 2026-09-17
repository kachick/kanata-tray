module github.com/rszyma/kanata-tray

go 1.25.0

// Using github.com/gogpu/systray for pure Go / zero-CGO systray support across platforms.
// Previously used github.com/getlantern/systray (unmaintained, required CGO on Linux/macOS)
// and fyne.io/systray (had ~1s icon switching delay on Linux).

require (
	github.com/elliotchance/orderedmap/v2 v2.2.0
	github.com/expr-lang/expr v1.17.8
	github.com/go-chi/chi/v5 v5.2.1
	github.com/gogpu/systray v0.3.0
	github.com/k0kubun/pp/v3 v3.2.0
	github.com/kirsle/configdir v0.0.0-20170128060238-e45d2f54772f
	github.com/kr/pretty v0.3.1
	github.com/labstack/gommon v0.4.2
	github.com/mattn/go-isatty v0.0.20
	github.com/pelletier/go-toml/v2 v2.4.0
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966
	github.com/spf13/pflag v1.0.6
)

require (
	github.com/go-webgpu/goffi v0.6.3 // indirect
	github.com/godbus/dbus/v5 v5.2.2 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
