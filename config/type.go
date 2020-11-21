package config

import (
	"errors"

	"github.com/mebiusashan/beaker/common"
)

const checkDatabase = 0x1
const checkWebsite = 0x2
const checkAuthInfo = 0x4
const checkServer = 0x8
const checkRedis = 0x10

type ConfigData struct {
	Database Database
	Server   Server
	Redis    Redis
	Website  Website
	AuthInfo Auth
}

type Database struct {
	DB_URL       string
	DB_USER      string
	DB_PW        string
	DB_NAME      string
	MAX_IDLE_NUM int
	MAX_OPEN_NUM int
}

type Server struct {
	PORT string
	URL  string
}

type Redis struct {
	REDIS_IP     string
	REDIS_PORT   string
	REDIS_PREFIX string
	EXPIRE_TIME  int
}

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

type Auth struct {
	Name         string
	Password     string
	ServerKeyDir string
	ClientKeyDir string
	ConfigPath   string
}

type BaseConfig interface {
	check() error
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
		t.MAX_IDLE_NUM = common.Def_Database_MAX_IDLE_NUM
	}
	if t.MAX_OPEN_NUM <= 0 {
		t.MAX_OPEN_NUM = common.Def_Database_MAX_OPEN_NUM
	}
	return nil
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
		t.EXPIRE_TIME = common.Def_Redis_EXPIRE_TIME
	}
	return nil
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
		t.SITE_NAME = common.Def_Website_SITE_NAME
	}
	if t.SITE_DES == "" {
		t.SITE_DES = common.Def_Website_SITE_DES
	}
	if t.SITE_FOOTER == "" {
		t.SITE_FOOTER = common.Def_Website_SITE_FOOTER
	}
	if t.INDEX_LIST_NUM <= 0 {
		t.INDEX_LIST_NUM = common.Def_Website_INDEX_LIST_NUM
	}
	if t.TWEET_NUM_ONE_PAGE <= 0 {
		t.TWEET_NUM_ONE_PAGE = common.Def_Website_TWEET_NUM_ONE_PAGE
	}
	return nil
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
