package cert

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"math/rand"
)

func CreateDesKey() []byte {
	key := make([]byte, 24)
	rand.Read(key)
	return key
}

// 3DES加密
func TripleDesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:8])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// 3DES解密
func TripleDesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:8])
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func PKCS5UnPadding(src []byte) []byte {
	//length := len(origData)
	//// 去掉最后一个字节 unpadding 次
	//unpadding := int(origData[length-1])
	//fmt.Println("长度", length)
	//return origData[:(length - unpadding)]
	length := len(src)
	//fmt.Println("测试长度",length, src[length-1])
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
