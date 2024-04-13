package flock

import (
	"fmt"
	"os"
)

var pidfile *Flock

// PidLock 锁定pid文件后写入pid信息
func PidLock(path string) error {
	if path == "" {
		return nil
	}

	pid := os.Getpid()
	info, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) { // 文件不存在
			return err
		}
	} else if info.IsDir() {
		return fmt.Errorf("pid file path is folder. ")
	} else if !info.Mode().IsRegular() {
		return fmt.Errorf("pid file is not regular file. ")
	}
	lock := New(path)
	locked, _ := lock.TryLock()
	if !locked {
		return fmt.Errorf("lock pid file failed. ")
	}
	pidfile = lock
	// 锁定成功,写入pid
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("open pid file failed. ")
	}
	file.WriteString(fmt.Sprintf("%d", pid))
	file.Sync()
	file.Close()
	return nil
}

func PidUnLock() error {
	if pidfile == nil {
		return nil
	}
	err := pidfile.Unlock()
	if err != nil {
		return err
	}
	path := pidfile.Path()
	pidfile = nil
	return os.Remove(path)
}
