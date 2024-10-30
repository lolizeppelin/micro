package utils

import (
	"encoding/json"
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

func LoadJson(path string, payload any) error {
	if payload == nil {
		return fmt.Errorf("LoadJson: nil payload")
	}
	if ok, err := PathIsRegularFile(path); !ok {
		if err != nil {
			return fmt.Errorf("LoadJson: path %s is load failed: %w", path, err)
		}
		return fmt.Errorf("LoadJson: path %s is not a regular file", path)
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("LoadJson: path %s is open failed: %w", path, err)
	}
	defer file.Close()
	var buff []byte
	buff, err = io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("LoadJson: path %s is read failed: %w", path, err)
	}
	return json.Unmarshal(buff, payload)
}

func SaveJson(path string, payload any, overwrite ...bool) error {
	if payload == nil {
		return fmt.Errorf("SaveJson: nil payload")
	}
	buff, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("SaveJson: path %s is read failed: %w", path, err)
	}
	var info os.FileInfo
	info, err = PathFileInfo(path)
	if err != nil {
		return fmt.Errorf("SaveJson: path %s is check failed: %w", path, err)
	}
	if info != nil {
		if len(overwrite) > 0 && overwrite[0] {
			return fmt.Errorf("SaveJson: path %s exist", path)
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("SaveJson: path %s not regular file", path)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("SaveJson: path %s is open failed: %w", path, err)
	}
	defer file.Close()
	_, err = file.Write(buff)
	if err != nil {
		return fmt.Errorf("SaveJson: path %s is read failed: %w", path, err)
	}
	return nil
}
