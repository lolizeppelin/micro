package utils

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
)

// UnZip 解压到文件夹
func UnZip(buf *bytes.Reader, output string) error {
	//buf := bytes.NewReader(buff)
	unzip, err := zip.NewReader(buf, buf.Size())
	if err != nil {
		return err
	}
	for _, f := range unzip.File {
		path := filepath.Join(output, f.Name)
		if f.FileInfo().IsDir() {
			if err = os.MkdirAll(path, f.Mode()); err != nil {
				return err
			}
		} else {
			var r io.ReadCloser
			r, err = f.Open()
			if err != nil {
				return err
			}
			defer r.Close()

			var w *os.File
			w, err = os.OpenFile(
				path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer w.Close()

			_, err = io.Copy(w, r)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Zip(path, output string) error {

	//writer, err := os.Create(output)
	//if err != nil {
	//	return err
	//}
	//defer writer.Close()
	//
	//x := zip.NewWriter(writer)

	return nil

}
