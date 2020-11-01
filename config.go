// 配置库
package beaker

import (
	"errors"
	"fmt"

	"github.com/BurntSushi/toml"
)

const checkDatabase = 0x1
const checkWebsite = 0x2
const checkAuthInfo = 0x4
const checkServer = 0x8
const checkRedis = 0x10

type BaseConfig interface {
	check() error
}

//配置数据
type ConfigData struct {
	Database Database
	Server   Server
	Redis    Redis
	Website  Website
	AuthInfo Auth
}

//数据库配置
type Database struct {
	DB_URL       string
	DB_USER      string
	DB_PW        string
	DB_NAME      string
	MAX_IDLE_NUM int
	MAX_OPEN_NUM int
}

func (t *Database) check() error {
	if t.DB_URL == "" {
		return errors.New("Database's DB_URL is empty.")
	}
	if t.DB_USER == "" {
		return errors.New("Database's DB_USER is empty.")
	}
	if t.DB_PW == "" {
		return errors.New("Database's DB_PW is empty.")
	}
	if t.DB_NAME == "" {
		return errors.New("Database's DB_NAME is empty.")
	}
	if t.MAX_IDLE_NUM <= 0 {
		t.MAX_IDLE_NUM = def_Database_MAX_IDLE_NUM
	}
	if t.MAX_OPEN_NUM <= 0 {
		t.MAX_OPEN_NUM = def_Database_MAX_OPEN_NUM
	}
	return nil
}

//服务器配置
type Server struct {
	PORT string
	URL  string
}

func (t *Server) check() error {
	if t.PORT == "" {
		return errors.New("Server's PORT is empty.")
	}
	if t.URL == "" {
		return errors.New("Server's URL is empty.")
	}
	return nil
}

//Redis配置
type Redis struct {
	REDIS_IP     string
	REDIS_PORT   string
	REDIS_PREFIX string
	EXPIRE_TIME  int
}

func (t *Redis) check() error {
	if t.REDIS_IP == "" {
		return errors.New("Redis's REDIS_IP is empty.")
	}
	if t.REDIS_PORT == "" {
		return errors.New("Redis's REDIS_PORT is empty.")
	}
	if t.REDIS_PREFIX == "" {
		return errors.New("Redis's REDIS_PREFIX is empty.")
	}
	if t.EXPIRE_TIME <= 0 {
		t.EXPIRE_TIME = def_Redis_EXPIRE_TIME
	}
	return nil
}

//网站信息
type Website struct {
	SITE_NAME          string
	SITE_URL           string
	SITE_DES           string
	SITE_FOOTER        string
	INDEX_LIST_NUM     uint
	TEMP_FOLDER        string
	STATIC_FILE_FOLDER string
	TWEET_NUM_ONE_PAGE uint
	SITE_KEYWORDS      string
}

func (t *Website) check() error {
	if t.STATIC_FILE_FOLDER == "" {
		return errors.New("Website's STATIC_FILE_FOLDER is empty.")
	}
	if t.TEMP_FOLDER == "" {
		return errors.New("Website's TEMP_FOLDER is empty.")
	}
	if t.SITE_URL == "" {
		return errors.New("Website's SITE_URL is empty.")
	}
	if t.SITE_NAME == "" {
		t.SITE_NAME = def_Website_SITE_NAME
	}
	if t.SITE_DES == "" {
		t.SITE_DES = def_Website_SITE_DES
	}
	if t.SITE_FOOTER == "" {
		t.SITE_FOOTER = def_Website_SITE_FOOTER
	}
	if t.INDEX_LIST_NUM <= 0 {
		t.INDEX_LIST_NUM = def_Website_INDEX_LIST_NUM
	}
	if t.TWEET_NUM_ONE_PAGE <= 0 {
		t.TWEET_NUM_ONE_PAGE = def_Website_TWEET_NUM_ONE_PAGE
	}
	return nil
}

//管理用户信息
type Auth struct {
	Name         string
	Password     string
	ServerKeyDir string
	ClientKeyDir string
	ConfigPath   string
}

func (t *Auth) check() error {
	if t.Name == "" {
		return errors.New("Auth's Name is empty.")
	}
	if t.Password == "" {
		return errors.New("Auth's Password is empty.")
	}
	if t.ServerKeyDir == "" {
		return errors.New("Auth's ServerKeyDir is empty.")
	}
	if t.ClientKeyDir == "" {
		return errors.New("Auth's ClientKeyDir is empty.")
	}
	if t.ConfigPath == "" {
		return errors.New("Auth's ConfigPath is empty.")
	}
	return nil
}

//创建一个新的配置对象，需要配置文件路径
func NewWithPath(path string, check byte) (*ConfigData, error) {
	var config ConfigData
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		return nil, errors.New("Config File Load Failed: " + err.Error())
	}

	if checkDatabase&check == checkDatabase {
		if err = checkWithDefConfig(&config.Database); err != nil {
			return nil, err
		}
	}
	if checkWebsite&check == checkWebsite {
		if err = checkWithDefConfig(&config.Website); err != nil {
			return nil, err
		}
	}
	if checkAuthInfo&check == checkAuthInfo {
		fmt.Println("checkAuthInfo")
		if err = checkWithDefConfig(&config.AuthInfo); err != nil {
			return nil, err
		}
	}
	if checkServer&check == checkServer {
		if err = checkWithDefConfig(&config.Server); err != nil {
			return nil, err
		}
	}
	if checkRedis&check == checkRedis {
		if err = checkWithDefConfig(&config.Redis); err != nil {
			return nil, err
		}
	}

	return &config, nil
}

func checkWithDefConfig(c BaseConfig) error {
	if c == nil {
		return nil
	}
	return c.check()
}
