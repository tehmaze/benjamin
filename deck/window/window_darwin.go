//go:build darwin

package window

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
)

func (Window) Manufacturer() string { return "Apple" }
func (Window) Product() string      { return "macOS" }

func (Window) SerialNumber() string {
	buf := &bytes.Buffer{}
	err := run(buf, os.Stderr, "ioreg", "-rd1", "-c", "IOPlatformExpertDevice")
	if err != nil {
		return ""
	}
	return extractID(buf.String())
}

func extractID(lines string) string {
	for _, line := range strings.Split(lines, "\n") {
		if strings.Contains(line, "IOPlatformSerialNumber") {
			parts := strings.SplitAfter(line, `" = "`)
			if len(parts) == 2 {
				return strings.TrimRight(parts[1], `"`)
			}
		}
	}
	return ""
}

// run wraps `exec.Command` with easy access to stdout and stderr.
func run(stdout, stderr io.Writer, cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stdin = os.Stdin
	c.Stdout = stdout
	c.Stderr = stderr
	return c.Run()
}
