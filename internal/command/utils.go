package command

import (
	"bytes"
	"io"
	"path"
	"regexp"
	"runtime"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/mebiusashan/beaker/internal/common"
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

func convMarkdownImage(markdown []byte, mdPath string) (string, []common.ImgInfo) {
	imgPaths := FindImageURL(markdown)
	mdStr := string(markdown)
	filenameWithSuffix := path.Base(mdPath)
	for i := 0; i < len(imgPaths); i++ {
		imgPaths[i].Read(mdPath[0 : len(mdPath)-len(filenameWithSuffix)])
		if imgPaths[i].Readed {
			mdStr = strings.ReplaceAll(mdStr, imgPaths[i].Path, imgPaths[i].Md5+imgPaths[i].Suffix)
		}
	}

	if runtime.GOOS == "windows" {
		lfReader := DOS2Unix(strings.NewReader(mdStr))
		buf := new(bytes.Buffer)
		buf.ReadFrom(lfReader)
		mdStr = buf.String()
	}

	return mdStr, imgPaths
}

type byteReader struct {
	io.Reader
	buf [1]byte
}

func (b *byteReader) ReadByte() (byte, error) {
	_, err := io.ReadFull(b.Reader, b.buf[:])
	return b.buf[0], err
}

type dos2unix struct {
	r    io.ByteReader
	b    bool
	char byte
}

func (d *dos2unix) Read(b []byte) (int, error) {
	var n int
	for len(b) > 0 {
		if d.b {
			b[0] = d.char
			d.b = false
			b = b[1:]
			n++
			continue
		}
		c, err := d.r.ReadByte()
		if err != nil {
			return n, err
		}
		if c == '\r' {
			d.char, err = d.r.ReadByte()
			if err != io.EOF {
				if err != nil {
					return n, err
				}
				if d.char == '\n' {
					c = '\n'
				} else {
					d.b = true
				}
			}
		}
		b[0] = c
		b = b[1:]
		n++
	}
	return n, nil
}

// DOS2Unix wraps a byte reader with a reader that replaces all instances of
// \r\n with \n
func DOS2Unix(r io.Reader) io.Reader {
	br, ok := r.(io.ByteReader)
	if !ok {
		br = &byteReader{Reader: r}
	}
	return &dos2unix{r: br}
}
