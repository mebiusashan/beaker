package cmd

import (
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/mebiusashan/beaker/common"
)

func FindImageURL(markdown []byte) []common.ImgInfo {
	var imgurls []common.ImgInfo
	var imgflag = regexp.MustCompile(`!\[[^\]]*\]\((.*?)\s*("(?:.*[^"])")?\s*\)`)
	rel := imgflag.FindAll(markdown, -1)
	for i := 0; i < len(rel); i++ {
		str := string(rel[i])
		r := strings.Index(str, "](")
		if r != -1 {
			path := str[r+2 : len(str)-1]
			if !govalidator.IsRequestURL(path) && !govalidator.IsIP(path) {
				var info = common.ImgInfo{Path: path, Md5: "", Base64: "", Readed: false}
				imgurls = append(imgurls, info)
			}
		}
	}

	return imgurls
}
