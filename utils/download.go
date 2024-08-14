package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

func DownloadToBuffer(url string, writer io.Writer) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	// Writer the body to file
	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

// DownloadToFile 下载文件
func DownloadToFile(url, path string) (err error) {
	var info os.FileInfo
	if info, err = PathFileInfo(path); err != nil || info != nil {
		err = fmt.Errorf("out path error")
		return
	}
	var f *os.File
	f, err = os.Create(path)
	if err != nil {
		return err
	}

	defer func() {
		f.Close()
		if err != nil {
			os.Remove(path)
		}
	}()

	var resp *http.Response
	resp, err = http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Check server response
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", resp.Status)
		return
	}
	// Writer the body to file
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return
	}
	return
}

// DownloadToFolder 下载文件到目录
func DownloadToFolder(url, folder string) (err error) {
	if ok, _ := PathIsDir(folder); !ok {
		return fmt.Errorf("output path is not folder")
	}

	var resp *http.Response
	resp, err = http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// Check server response
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", resp.Status)
		return
	}
	name := resp.Header.Get("Content-Disposition")
	if name == "" {
		name = path.Base(url)
	}
	if !IsSafePathName(name) {
		err = fmt.Errorf("illegal file name %s", name)
		return
	}
	// 检查文件名是否包含相对目录
	output := path.Join(folder, name)
	var info os.FileInfo
	if info, err = PathFileInfo(output); err != nil || info != nil {
		err = fmt.Errorf("out path error")
		return
	}

	var f *os.File
	f, err = os.Create(output)
	if err != nil {
		return
	}

	defer func() {
		f.Close()
		if err != nil {
			os.Remove(output)
		}
	}()

	// Writer the body to file
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return
	}
	return
}
