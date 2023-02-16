//go:build linux

package window

import (
	"io"
	"strings"
)

func (Window) Manufacturer() string { return "GNU" }
func (Window) Product() string      { return "Linux" }

func (Window) SerialNumber() string {
	for _, name := range []string{
		"/etc/machine-id",          // systemd
		"/var/lib/dbus/machine-id", // dbus
	} {
		if n, _ := io.ReadFile(name); n != "" {
			return strings.TrimSpace(n)
		}
	}
	return ""
}
