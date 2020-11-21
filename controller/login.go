package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/cert"
	"github.com/mebiusashan/beaker/common"
)

type LoginController struct {
	BaseController
}

var loginKey string = ""

func (ct *LoginController) Ping(c *gin.Context) {
	pubKey, err := ioutil.ReadFile(ct.Context.Config.AuthInfo.ServerKeyDir + common.SERVER_PUBLIC_KEY)
	if err != nil {
		writeFail(c, err.Error())
		return
	}
	writeSucc(c, "", cert.Base64Encode(pubKey))
}

func (ct *LoginController) Login(c *gin.Context) {
	pri, err := ioutil.ReadFile(ct.Context.Config.AuthInfo.ServerKeyDir + common.SERVER_PRIVATE_KEY)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	data, err = cert.Base64Decode(string(data))
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	var jsonData common.LoginReq
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	if jsonData.UN == ct.Context.Config.AuthInfo.Name {
		pw := cert.MD5([]byte(ct.Context.Config.AuthInfo.Password))
		if pw == jsonData.PW {
			//账号正确
			//从redis里度des key，没有就创建
			//用dk，加密 des key，发回去
			desKey := loginKey
			if desKey == "" {
				key := cert.CreateDesKey()
				desKey = base64.StdEncoding.EncodeToString(key)
				loginKey = desKey
			}

			key, _ := base64.StdEncoding.DecodeString(desKey)
			clientDesKey64, err := cert.Base64Decode(jsonData.DK)
			if err != nil {
				writeFail(c, "Decoding failed"+err.Error())
				return
			}

			clientDesKey, err := cert.RSADecrypt(pri, []byte(clientDesKey64))
			if err != nil {
				writeFail(c, "Decoding failed"+err.Error())
				return
			}

			serverDesKeyM, err := cert.TripleDesEncrypt(key, clientDesKey)
			if err != nil {
				writeFail(c, "Decoding failed"+err.Error())
				return
			}

			serverDesKey64 := cert.Base64Encode(serverDesKeyM)
			writeSucc(c, "", serverDesKey64)
		}
	}
	if err != nil {
		writeFail(c, "error")
		return
	}
}

func (ct *LoginController) Check(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	desKey := loginKey
	if desKey == "" {
		writeFail(c, "Need Login")
		return
	}

	key, err := base64.StdEncoding.DecodeString(desKey)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	data64, err := cert.Base64Decode(string(data))
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	sl, err := cert.TripleDesDecrypt(data64, key)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	sl123 := string(sl) + "123"
	rel, err := cert.TripleDesEncrypt([]byte(sl123), key)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	writeSucc(c, "", cert.Base64Encode(rel))
}
