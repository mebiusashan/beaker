package controller

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/internal/cert"
	"github.com/mebiusashan/beaker/internal/common"
)

type LoginController struct {
	BaseController
}

type LoginInfo struct {
	loginKey   string
	createTime int64
}

var sessions = struct {
	sync.RWMutex
	values map[string]LoginInfo
}{values: make(map[string]LoginInfo)}

func newSession(key string) string {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return ""
	}
	token := hex.EncodeToString(buf)
	sessions.Lock()
	sessions.values[token] = LoginInfo{loginKey: key, createTime: time.Now().Unix()}
	sessions.Unlock()
	return token
}

func sessionKey(token string, expire int64) (string, bool) {
	if token == "" {
		return "", false
	}
	sessions.RLock()
	info, ok := sessions.values[token]
	sessions.RUnlock()
	if !ok || time.Now().Unix()-info.createTime >= expire {
		sessions.Lock()
		delete(sessions.values, token)
		sessions.Unlock()
		return "", false
	}
	return info.loginKey, true
}

func (ct *LoginController) Ping(c *gin.Context) {
	pubKey, err := os.ReadFile(ct.Context.Config.AuthInfo.ServerKeyDir + common.SERVER_PUBLIC_KEY)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, common.ErrorCodeInternal, err)
		return
	}
	writeSucc(c, "", cert.Base64Encode(pubKey))
}

func (ct *LoginController) Login(c *gin.Context) {
	pri, err := os.ReadFile(ct.Context.Config.AuthInfo.ServerKeyDir + common.SERVER_PRIVATE_KEY)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, common.ErrorCodeInternal, err)
		return
	}

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		ErrorFromCode(c, http.StatusBadRequest, common.ErrorCodeInvalidRequest, err.Error())
		return
	}

	data, err = cert.Base64Decode(string(data))
	if err != nil {
		ErrorFromCode(c, http.StatusBadRequest, common.ErrorCodeDecode, "invalid login payload")
		return
	}

	var jsonData common.LoginReq
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		ErrorFromCode(c, http.StatusBadRequest, common.ErrorCodeInvalidRequest, err.Error())
		return
	}

	if jsonData.UN == ct.Context.Config.AuthInfo.Name {
		pw := cert.MD5([]byte(ct.Context.Config.AuthInfo.Password))
		if pw == jsonData.PW {
			key := cert.CreateDesKey()
			if len(key) == 0 {
				ErrorFromCode(c, http.StatusInternalServerError, common.ErrorCodeInternal, "failed to create session key")
				return
			}
			desKey := base64.StdEncoding.EncodeToString(key)
			clientDesKey64, err := cert.Base64Decode(jsonData.DK)
			if err != nil {
				ErrorFromCode(c, http.StatusBadRequest, common.ErrorCodeDecode, "Decoding failed")
				return
			}

			clientDesKey, err := cert.RSADecrypt(pri, []byte(clientDesKey64))
			if err != nil {
				ErrorFromCode(c, http.StatusBadRequest, common.ErrorCodeDecode, "Decoding failed")
				return
			}

			serverDesKeyM, err := cert.TripleDesEncrypt(key, clientDesKey)
			if err != nil {
				ErrorFromCode(c, http.StatusBadRequest, common.ErrorCodeDecode, "Decoding failed")
				return
			}

			serverDesKey64 := cert.Base64Encode(serverDesKeyM)
			resp := new(common.SuccMsgResp)
			resp.Code = common.SUCC
			resp.RequestID = requestID(c)
			resp.Data = serverDesKey64
			resp.SessionToken = newSession(desKey)
			if resp.SessionToken == "" {
				ErrorFromCode(c, http.StatusInternalServerError, common.ErrorCodeInternal, "failed to create session")
				return
			}
			c.JSON(http.StatusOK, resp)
			return
		}
	}
	ErrorFromCode(c, http.StatusUnauthorized, common.ErrorCodeUnauthorized, "invalid username or password")
}

func (ct *LoginController) Check(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		ErrorFromCode(c, http.StatusBadRequest, common.ErrorCodeInvalidRequest, err.Error())
		return
	}

	desKey, ok := sessionKey(c.GetHeader("X-Beaker-Session"), ct.Context.Config.AuthInfo.EXPIRE_TIME)
	if !ok {
		ErrorFromCode(c, http.StatusUnauthorized, common.ErrorCodeUnauthorized, "Need Login")
		return
	}

	key, err := base64.StdEncoding.DecodeString(desKey)
	if err != nil {
		ErrorFromCode(c, http.StatusBadRequest, common.ErrorCodeDecode, "No credit request")
		return
	}

	data64, err := cert.Base64Decode(string(data))
	if err != nil {
		ErrorFromCode(c, http.StatusBadRequest, common.ErrorCodeDecode, "No credit request")
		return
	}

	sl, err := cert.TripleDesDecrypt(data64, key)
	if err != nil {
		ErrorFromCode(c, http.StatusBadRequest, common.ErrorCodeDecode, "No credit request")
		return
	}

	sl123 := string(sl) + "123"
	rel, err := cert.TripleDesEncrypt([]byte(sl123), key)
	if err != nil {
		ErrorFromCode(c, http.StatusBadRequest, common.ErrorCodeDecode, "No credit request")
		return
	}

	writeSucc(c, "", cert.Base64Encode(rel))
}
