package cow

import (
	"fmt"
	"golang.org/x/sys/windows"
	"io"
	"os"
	"unsafe"
)

// ReFS block cloning requires cluster-aligned regions under 4 GiB. Clone the
// aligned body and copy the short tail. NTFS/other unsupported volumes fall back
// through the common copy path. Neither operation creates hard links.
func cloneFile(source, destination string) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	info, err := in.Stat()
	if err != nil {
		return err
	}
	path, err := windows.UTF16PtrFromString(source)
	if err != nil {
		return err
	}
	volume := make([]uint16, 32768)
	if err := windows.GetVolumePathName(path, &volume[0], uint32(len(volume))); err != nil {
		return err
	}
	var sectors, bytes, free, total uint32
	proc := windows.NewLazySystemDLL("kernel32.dll").NewProc("GetDiskFreeSpaceW")
	ok, _, callErr := proc.Call(uintptr(unsafe.Pointer(&volume[0])), uintptr(unsafe.Pointer(&sectors)), uintptr(unsafe.Pointer(&bytes)), uintptr(unsafe.Pointer(&free)), uintptr(unsafe.Pointer(&total)))
	if ok == 0 {
		return callErr
	}
	cluster := int64(sectors) * int64(bytes)
	if cluster <= 0 || info.Size() < cluster {
		return fmt.Errorf("file has no complete cloneable cluster")
	}
	out, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer out.Close()
	if err := out.Truncate(info.Size()); err != nil {
		return err
	}
	aligned := info.Size() / cluster * cluster
	chunk := int64(1<<30) / cluster * cluster
	if chunk == 0 {
		return fmt.Errorf("unsupported cluster size")
	}
	for offset := int64(0); offset < aligned; {
		length := min(chunk, aligned-offset)
		data := struct {
			Source                             windows.Handle
			SourceOffset, TargetOffset, Length int64
		}{windows.Handle(in.Fd()), offset, offset, length}
		var returned uint32
		if err := windows.DeviceIoControl(windows.Handle(out.Fd()), windows.FSCTL_DUPLICATE_EXTENTS_TO_FILE, (*byte)(unsafe.Pointer(&data)), uint32(unsafe.Sizeof(data)), nil, 0, &returned, nil); err != nil {
			return err
		}
		offset += length
	}
	if _, err := in.Seek(aligned, io.SeekStart); err != nil {
		return err
	}
	if _, err := out.Seek(aligned, io.SeekStart); err != nil {
		return err
	}
	_, err = io.Copy(out, in)
	return err
}
