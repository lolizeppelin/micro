package utils

import (
	"archive/zip"
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ZipCompress(data []byte) ([]byte, error) {
	// 创建一个字节缓冲区来存储压缩数据
	var buf bytes.Buffer
	// 创建一个 zlib.Writer
	zw := zlib.NewWriter(&buf)
	// 将数据写入 zlib.Writer
	_, err := zw.Write(data)
	if err != nil {
		return nil, err
	}
	// 关闭 zlib.Writer 以完成压缩
	err = zw.Close()
	if err != nil {
		return nil, err
	}
	// 返回压缩后的字节切片
	return buf.Bytes(), nil
}

func ZipDecompress(data []byte) ([]byte, error) {
	buf := bytes.NewReader(data)
	zr, err := zlib.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	var out bytes.Buffer
	_, err = io.Copy(&out, zr)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// UnZip 解压到文件夹
func UnZip(source, output string) error {
	file, err := os.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return err
	}
	unzip, err := zip.NewReader(file, info.Size())
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

func ZipFileBuff(name string, data []byte) ([]byte, error) {
	buff := new(bytes.Buffer)
	// Create a new zip writer
	zw := zip.NewWriter(buff)
	// Create a new zip entry
	entity, err := zw.Create(name)
	if err != nil {
		zw.Close()
		return nil, fmt.Errorf("failed to create zip entry: %w", err)
	}
	_, err = entity.Write(data)
	if err != nil {
		zw.Close()
		return nil, fmt.Errorf("failed to write to zip entry: %w", err)
	}
	err = zw.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close zip writer: %w", err)
	}
	// Copy the data from buf to the zip entry
	return buff.Bytes(), nil
}

// Zip 压缩指定文件、文件夹
func Zip(source, output string) error {
	var fileInfo os.FileInfo
	if info, err := PathFileInfo(source); err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	} else {
		fileInfo = info
	}

	var out *os.File
	if file, err := os.Create(output); err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	} else {
		out = file
	}
	defer out.Close()

	if fileInfo.IsDir() {
		// path zip
		root := filepath.Base(source)
		// Create a new zip writer
		zw := zip.NewWriter(out)
		defer zw.Close()
		// Walk through the file or directory
		return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Create a zip header
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			// Set the header name
			if root != "" {
				header.Name = filepath.Join(root, strings.TrimPrefix(path, source))
			} else {
				header.Name = filepath.Base(source)
			}

			// If the current file is a directory, add a trailing slash to the name
			if info.IsDir() {
				header.Name += "/"
			} else {
				// Set the compression method for files
				header.Method = zip.Deflate
			}

			// Create a writer for the current file in the zip
			writer, err := zw.CreateHeader(header)
			if err != nil {
				return err
			}

			var file *os.File
			// If the current file is not a directory, write its content to the zip
			if !info.IsDir() {
				file, err = os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()

				_, err = io.Copy(writer, file)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	// Compress file
	file, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer file.Close()

	zw := zip.NewWriter(out)
	defer zw.Close()
	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}
	header.Name = filepath.Base(source)
	writer, err := zw.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, file)
	return err
}
