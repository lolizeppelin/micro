package utils

import (
	"fmt"
	"io"
	"os"
)

func CopyFile(src, dst string) (err error) {
	if _, err = os.Stat(dst); err == nil {
		return fmt.Errorf("CopyFile: file exist")
	} else if !os.IsNotExist(err) {
		return
	}
	var info os.FileInfo
	info, err = os.Stat(src)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", info.Name(), info.Mode().String())
	}

	var in *os.File
	in, err = os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	var out *os.File
	out, err = os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		cErr := out.Close()
		if err == nil {
			err = cErr
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return nil
}
