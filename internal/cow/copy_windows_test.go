package cow

import (
	"testing"
	"unsafe"

	"golang.org/x/sys/windows"
)

func init() {
	volumeCluster = func(t *testing.T, dir string) int64 {
		t.Helper()
		path, err := windows.UTF16PtrFromString(dir)
		if err != nil {
			t.Fatal(err)
		}
		volume := make([]uint16, 32768)
		if err := windows.GetVolumePathName(path, &volume[0], uint32(len(volume))); err != nil {
			t.Fatal(err)
		}
		var sectors, bytes, free, total uint32
		proc := windows.NewLazySystemDLL("kernel32.dll").NewProc("GetDiskFreeSpaceW")
		ok, _, callErr := proc.Call(uintptr(unsafe.Pointer(&volume[0])), uintptr(unsafe.Pointer(&sectors)), uintptr(unsafe.Pointer(&bytes)), uintptr(unsafe.Pointer(&free)), uintptr(unsafe.Pointer(&total)))
		if ok == 0 {
			t.Fatal("GetDiskFreeSpaceW:", callErr)
		}
		cluster := int64(sectors) * int64(bytes)
		if cluster <= 0 {
			t.Fatalf("volume of %s reports no cluster size", dir)
		}
		return cluster
	}
}
