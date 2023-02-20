//go:build windows

package window

import (
	"golang.org/x/sys/windows/registry"
)

func (Window) Manufacturer() string { return "Microsoft" }
func (Window) Product() string      { return "Windows" }

func (Window) SerialNumber() string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Cryptography`, registry.QUERY_VALUE|registry.WOW64_64KEY)
	if err != nil {
		return ""
	}
	defer k.Close()

	s, _, _ := k.GetStringValue("MachineGuid")
	return s
}
