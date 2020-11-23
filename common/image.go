package common

import (
	"bufio"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"io"
	"os"
	"path"
)

type ImgInfo struct {
	Path   string
	Suffix string
	Md5    string
	Base64 string
	Readed bool
}

func (info *ImgInfo) Read(parentPath string) {
	filenameWithSuffix := path.Base(info.Path)
	info.Suffix = path.Ext(filenameWithSuffix)
	filePath := parentPath + "/" + info.Path
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
	info.Base64 = base64.StdEncoding.EncodeToString(buf)

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return
	}
	hashInBytes := hash.Sum(nil)[:16]
	info.Md5 = hex.EncodeToString(hashInBytes)
	info.Readed = true
}
