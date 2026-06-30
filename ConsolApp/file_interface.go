package main

import (
	"os"
	"strings"
)

type FileToSend struct {
	totalSize int64
	fileName  string
}

func IsPathValid(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

func GetFileMetadata(filepath string) *FileToSend {
	splitStr := strings.Split("/", filepath)
	fileName := splitStr[len(splitStr)-1]
	info, _ := os.Stat(filepath)
	var fileSize int64 = info.Size()
	return &FileToSend{
		totalSize: fileSize,
		fileName:  fileName,
	}

}

func sendFile(filepath string) bool {
	fileMetadata := GetFileMetadata(filepath)

	return true
}
