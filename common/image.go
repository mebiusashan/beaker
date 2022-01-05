package common

import (
	"bufio"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"io"
	"os"
	"path"
	"path/filepath"
)

type ImgInfo struct {
	Path   string
	Suffix string
	Md5    string
	Base64 string
	Readed bool
}

func (imgInfo *ImgInfo) Read(parentPath string) {
	filenameWithSuffix := path.Base(imgInfo.Path)
	imgInfo.Suffix = path.Ext(filenameWithSuffix)
	filePath := filepath.Join(parentPath, imgInfo.Path)
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	fInfo, err := file.Stat()
	if err != nil {
		return
	}
	var size int64 = fInfo.Size()
	buf := make([]byte, size)
	fReader := bufio.NewReader(file)
	_, err = fReader.Read(buf)
	if err != nil {
		return
	}
	imgInfo.Base64 = base64.StdEncoding.EncodeToString(buf)

	md5file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer md5file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, md5file); err != nil {
		return
	}
	hashInBytes := hash.Sum(nil)[:16]
	imgInfo.Md5 = hex.EncodeToString(hashInBytes)
	imgInfo.Readed = true
}
