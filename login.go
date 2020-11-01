package beaker

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/cert"
)

/**
ping: c -> 请求s公钥-> s
ping: c <-   s公钥  <- s
login: c -> 用s公钥加密c私钥和用户名密码 -> s
login: c <- 用c私钥加密用户名和私有公钥 <- s
check: c -> 用c公钥解密用户名验证，s私有公钥加密验证结果 ->s
check: c <- 用c私有私钥加密用户名 <- s

以后通信都是有 对应用户的私钥公钥
*/

const SERVER_PUBLIC_KEY = "pub.pem"
const SERVER_PRIVATE_KEY = "pri.pem"

func (ct *LoginCtrl) Ping(c *gin.Context) {
	fmt.Println(ct.ctrl.mvc.config.AuthInfo.ServerKeyDir + SERVER_PUBLIC_KEY)
	pubKey, err := ioutil.ReadFile(ct.ctrl.mvc.config.AuthInfo.ServerKeyDir + SERVER_PUBLIC_KEY)
	if err != nil {
		writeFail(c, err.Error())
		return
	}
	writeSucc(c, "", cert.Base64Encode(pubKey))
}

type LoginReq struct {
	DK string
	UN string
	PW string
}

func (ct *LoginCtrl) Login(c *gin.Context) {
	pri, err := ioutil.ReadFile(ct.ctrl.mvc.config.AuthInfo.ServerKeyDir + SERVER_PRIVATE_KEY)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		writeFail(c, err.Error())
		return
	}
	fmt.Println("收到的数据", string(data))

	data, err = cert.Base64Decode(string(data))
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("收到的数据", string(data))

	var jsonData LoginReq
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	if jsonData.UN == ct.ctrl.mvc.config.AuthInfo.Name {
		pw := cert.MD5([]byte(ct.ctrl.mvc.config.AuthInfo.Password))
		if pw == jsonData.PW {

			//账号正确
			//从redis里度des key，没有就创建
			//用dk，加密 des key，发回去

			desKey, err := ct.ctrl.mvc.cache.GET("dk", "dk")
			if err != nil {
				if err.Error() == "redigo: nil returned" {
					key := cert.CreateDesKey()
					fmt.Println("长---度", len(key))
					desKey = base64.StdEncoding.EncodeToString(key)
					ct.ctrl.mvc.cache.SETNX("dk", "dk", desKey, ct.ctrl.mvc.config.Redis.EXPIRE_TIME)
				} else {
					writeFail(c, err.Error())
					return
				}
			}

			key, _ := base64.StdEncoding.DecodeString(desKey)
			fmt.Println("生成base64的key1111", desKey)

			clientDesKey64, err := cert.Base64Decode(jsonData.DK)
			if err != nil {
				writeFail(c, "解密失败0："+err.Error())
				return
			}

			clientDesKey, err := cert.RSADecrypt(pri, []byte(clientDesKey64))
			if err != nil {
				writeFail(c, "解密失败1："+err.Error())
				return
			}

			serverDesKeyM, err := cert.TripleDesEncrypt(key, clientDesKey)
			if err != nil {
				writeFail(c, "解密失败2："+err.Error())
				return
			}

			serverDesKey64 := cert.Base64Encode(serverDesKeyM)

			fmt.Println("发送的加密Key", serverDesKey64)
			writeSucc(c, "", serverDesKey64)
		}
	}
	if err != nil {
		writeFail(c, "error")
		return
	}
}

func (ct *LoginCtrl) Check(c *gin.Context) {

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	fmt.Println("111收到的body", string(data))

	desKey, err := ct.ctrl.mvc.cache.GET("dk", "dk")
	if err != nil {
		fmt.Println("法师打发第三方", err)
		writeFail(c, err.Error())
		return
	}

	fmt.Println("2222查到的key", desKey)

	key, err := base64.StdEncoding.DecodeString(desKey)
	if err != nil {
		fmt.Println("redis 读取base64错误", err)
		writeFail(c, err.Error())
		return
	}
	//fmt.Println("生成base64的key", desKey)
	//fmt.Println("33333生成base64的key", string(data))

	data64, err := cert.Base64Decode(string(data))
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	fmt.Println("data64:::", data64)

	//dd, _ := base64.StdEncoding.DecodeString(string(data))
	//fmt.Println("解码body", string(dd))

	sl, err := cert.TripleDesDecrypt(data64, key)
	if err != nil {
		fmt.Println("解密错误", err)
		writeFail(c, err.Error())
		return
	}
	fmt.Println("解密后", string(sl))

	sl123 := string(sl) + "123"

	rel, err := cert.TripleDesEncrypt([]byte(sl123), key)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	writeSucc(c, "", cert.Base64Encode(rel))
}

func CMDPing() []byte {
	resp, err := http.Post(HOST+"/user/ping", "", strings.NewReader(""))
	if err != nil {
		//fmt.Println("ping", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//fmt.Println("ping", err)
	}

	var jsonData SuccMsg
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		//fmt.Println("ping", err)
	}

	if jsonData.Code != SUCC {
		fmt.Println(jsonData.Msg)
		return nil
	}

	pubKeyStr := jsonData.Data.(string)
	pubkey, err := cert.Base64Decode(pubKeyStr)
	if err != nil {
		//fmt.Println("ping", err)
	}

	//fmt.Println(string(pubkey))

	return pubkey
}

func CMDLogin(pubKey []byte) []byte {

	var fusername string
	var fpassword string

	fmt.Printf("请输入用户名:")

	fmt.Scanln(&fusername)

	fmt.Printf("请输入密码:")
	fmt.Scanln(&fpassword)

	clientDesKey := cert.CreateDesKey()
	//fmt.Println(string(pubKey))
	clientDesKeyM, err := cert.RSAEncryp(pubKey, clientDesKey)
	if err != nil {
		//fmt.Println("login", err)
	}

	jsonSendData := LoginReq{DK: cert.Base64Encode(clientDesKeyM), UN: fusername, PW: cert.MD5([]byte(fpassword))}

	//fmt.Println(jsonSendData)
	jsonByte, err := json.Marshal(jsonSendData)
	if err != nil {
		//fmt.Println("login", err)
	}

	jsonStr := cert.Base64Encode(jsonByte)
	resp, err := http.Post(HOST+"/user/login", "", strings.NewReader(jsonStr))
	if err != nil {
		//fmt.Println("login", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//fmt.Println("login", err)
	}

	var jsonData SuccMsg
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		//fmt.Println("login", err)
	}

	if jsonData.Code != SUCC {
		fmt.Println(jsonData.Msg)
		return nil
	}

	//fmt.Println("解64:", jsonData.Data.(string))
	serverDesKeyM, err := cert.Base64Decode(jsonData.Data.(string))
	if err != nil {
		//fmt.Println("login", err)
	}

	serverDesKey, err := cert.TripleDesDecrypt(serverDesKeyM, clientDesKey)
	if err != nil {
		//fmt.Println("login", err)
	}

	//fmt.Println("login:::", cert.Base64Encode(serverDesKey))

	return serverDesKey
}

func CMDCheck(serverDesKey []byte) bool {
	T := []byte("HackerBlog------====------")
	//fmt.Println("加密用的key", cert.Base64Encode(serverDesKey))
	desT, err := cert.TripleDesEncrypt(T, serverDesKey)
	if err != nil {
		//fmt.Println("check", err)
	}
	//fmt.Println("发送的机密后数据", desT)

	resp, err := http.Post(HOST+"/user/check", "", strings.NewReader(cert.Base64Encode(desT)))
	if err != nil {
		//fmt.Println("check", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//fmt.Println("check", err)
	}

	var jsonData SuccMsg
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		//fmt.Println("check", err)
	}

	if jsonData.Code != SUCC {
		//fmt.Println(jsonData.Msg)
		return false
	}
	//fmt.Println(jsonData)

	rel64, err := cert.Base64Decode(jsonData.Data.(string))
	if err != nil {
		//fmt.Println("check", err)
	}

	checkT, err := cert.TripleDesDecrypt(rel64, serverDesKey)
	if err != nil {
		//fmt.Println("check", err)
	}

	if string(checkT) == (string(T) + "123") {
		//将des key写入文件
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			//fmt.Println("check", err)
		}
		//fmt.Println(dir)
		path := dir + "/.hbkey"
		err = ioutil.WriteFile(path, serverDesKey, 0666)
		if err != nil {
			//fmt.Println("check", err)
		}

		return true
	}
	return false
}

func CMDGetLocalKey() ([]byte, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //exec.LookPath(os.Args[0])
	if err != nil {
		return nil, err
	}
	path := dir + "/.hbkey"
	key, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return key, nil
}
