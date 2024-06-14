package utils

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"time"
)

// FillRandomBytes 填充随机字节
func FillRandomBytes(b []byte) error {
	// 通过 crypto/rand 包填充随机字节
	_, err := rand.Read(b)
	return err
}

// ValidateC2 验证 C2 是否与 S1 匹配
func ValidateC2(s1, c2 []byte) bool {
	// 简单验证：C2 前 4 字节应与 S1 前 4 字节相同
	return binary.BigEndian.Uint32(s1) == binary.BigEndian.Uint32(c2)
}

// ReadFullWithTimeout 带超时的读取
func ReadFullWithTimeout(r io.Reader, buf []byte, timeout time.Duration) error {
	done := make(chan error, 1)
	go func() {
		_, err := io.ReadFull(r, buf)
		done <- err
	}()
	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		return errors.New("read timeout")
	}
}

// WriteWithTimeout 带超时的写入
func WriteWithTimeout(w io.Writer, buf []byte, timeout time.Duration) error {
	done := make(chan error, 1)
	go func() {
		_, err := w.Write(buf)
		done <- err
	}()
	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		return errors.New("write timeout")
	}
}

// LogError 记录错误日志
func LogError(err error) {
	if err != nil {
		log.Printf("error: %v", err)
	}
}
