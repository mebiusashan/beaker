package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/mebiusashan/beaker/common"
)

const def_RSA_LEN int = 2048

func RSAEncryp(pubKey []byte, origData []byte) ([]byte, error) {
	block, _ := pem.Decode(pubKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

func RSADecrypt(priKey []byte, ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode(priKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

func CheckRSAKey(pub []byte, pri []byte) {
	testStr := "test"
	rel, err := RSAEncryp(pub, []byte(testStr))
	common.Assert(err)
	rel, err = RSADecrypt(pri, rel)
	common.Assert(err)
	if string(rel) != testStr {
		common.Err("Secret key verification failed")
	}
}

func CreateRSAKeys() (pub []byte, pri []byte, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, def_RSA_LEN)
	if err != nil {
		return nil, nil, err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	priKey := pem.EncodeToMemory(block)

	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, nil, err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	pubKey := pem.EncodeToMemory(block)
	return pubKey, priKey, nil
}
