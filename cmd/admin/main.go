package main

import (
	"io/ioutil"
	"os"

	"github.com/mebiusashan/beaker/internal/cache"
	"github.com/mebiusashan/beaker/internal/cert"
	"github.com/mebiusashan/beaker/internal/common"
	"github.com/mebiusashan/beaker/internal/config"
	"github.com/mebiusashan/beaker/internal/controller"
	"github.com/mebiusashan/beaker/internal/database"
	"github.com/mebiusashan/beaker/internal/net"

	"github.com/gin-gonic/gin"
)

func main() {
	RunAdmin(true)
}

func RunAdmin(isRelease bool) {
	path := os.Getenv(common.ADMIN_ENV_KEY)
	config, err := config.NewWithPath(path, 0x1D)
	common.Assert(err)

	cac := cache.NewCache(config.Redis.REDIS_IP, config.Redis.REDIS_PORT, config.Redis.REDIS_PREFIX)
	model, err := database.NewDao(config.Database.DB_URL, config.Database.DB_USER, config.Database.DB_PW, config.Database.DB_NAME)
	common.Assert(err)
	model.SetMaxIdleConns(config.Database.MAX_IDLE_NUM)
	model.SetMaxOpenConns(config.Database.MAX_OPEN_NUM)

	context := controller.NewContext()
	context.Cache = cac
	context.Config = config
	context.Model = model

	pubHas, err := controller.PathExists(config.AuthInfo.ServerKeyDir + common.SERVER_PUBLIC_KEY)
	common.Assert(err)

	priHas, err := controller.PathExists(config.AuthInfo.ServerKeyDir + common.SERVER_PRIVATE_KEY)
	common.Assert(err)

	if !pubHas || !priHas {
		pub, pri, err := cert.CreateRSAKeys()
		common.Assert(err)
		err = ioutil.WriteFile(config.AuthInfo.ServerKeyDir+common.SERVER_PUBLIC_KEY, pub, 0666)
		common.Assert(err)
		err = ioutil.WriteFile(config.AuthInfo.ServerKeyDir+common.SERVER_PRIVATE_KEY, pri, 0666)
		common.Assert(err)
	}

	pub, err := ioutil.ReadFile(config.AuthInfo.ServerKeyDir + common.SERVER_PUBLIC_KEY)
	common.Assert(err)

	pri, err := ioutil.ReadFile(config.AuthInfo.ServerKeyDir + common.SERVER_PRIVATE_KEY)
	common.Assert(err)

	rel := cert.CheckRSAKey(pub, pri)
	if !rel {
		common.Err("Secret key verification failed")
	}

	if isRelease {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	//登录
	user := router.Group(net.ADMIN_GROUP_USER)
	user.POST(net.ADMIN_PING, context.Ctrl.LoginC.Ping)
	user.POST(net.ADMIN_LOGIN, context.Ctrl.LoginC.Login)
	user.POST(net.ADMIN_CHECK, context.Ctrl.LoginC.Check)

	adminr := router.Group(net.ADMIN_GROUP_ADMIN)
	adminr.Use(controller.LoginExpiredCheck())
	adminr.Use(controller.DecodeForAdmin())
	adminr.Use(controller.RefreshCache())

	adminr.POST(net.ADMIN_ART_ADD, context.Ctrl.ArtC.Add)
	adminr.POST(net.ADMIN_ART_RM, context.Ctrl.ArtC.Del)
	adminr.POST(net.ADMIN_ART_LIST, context.Ctrl.ArtC.All)
	adminr.POST(net.ADMIN_ART_DOWNLOAD, context.Ctrl.ArtC.Down)
	adminr.POST(net.ADMIN_ART_MODIFY, context.Ctrl.ArtC.Modify)

	adminr.POST(net.ADMIN_PAGE_ADD, context.Ctrl.PagC.Add)
	adminr.POST(net.ADMIN_PAGE_RM, context.Ctrl.PagC.Del)
	adminr.POST(net.ADMIN_PAGE_LIST, context.Ctrl.PagC.List)
	adminr.POST(net.ADMIN_PAGE_DOWNLOAD, context.Ctrl.PagC.Down)
	adminr.POST(net.ADMIN_PAGE_MODIFY, context.Ctrl.PagC.Modify)

	adminr.POST(net.ADMIN_TWEET_ADD, context.Ctrl.TweC.Add)
	adminr.POST(net.ADMIN_TWEET_RM, context.Ctrl.TweC.Del)
	adminr.POST(net.ADMIN_TWEET_LIST, context.Ctrl.TweC.List)

	adminr.POST(net.ADMIN_CAT_ADD, context.Ctrl.CatC.Add)
	adminr.POST(net.ADMIN_CAT_RM, context.Ctrl.CatC.Del)
	adminr.POST(net.ADMIN_CAT_LIST, context.Ctrl.CatC.All)
	adminr.POST(net.ADMIN_CAT_MODIFY, context.Ctrl.CatC.Update)

	adminr.POST(net.ADMIN_OPT, context.Ctrl.OptC.Info)
	adminr.POST(net.ADMIN_CLEAN, context.Ctrl.OptC.ClearCache)

	router.Run(config.Server.PORT)
}
