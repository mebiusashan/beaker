package main

import (
	"os"

	"github.com/mebiusashan/beaker/internal/cache"
	"github.com/mebiusashan/beaker/internal/common"
	"github.com/mebiusashan/beaker/internal/config"
	"github.com/mebiusashan/beaker/internal/controller"
	"github.com/mebiusashan/beaker/internal/database"
	"github.com/mebiusashan/beaker/internal/net"
	"github.com/mebiusashan/beaker/internal/view"

	"github.com/gin-gonic/gin"
)

func main() {
	RunServer(true)
}

func RunServer(isRelease bool) {
	path := os.Getenv(common.SERVER_ENV_KEY)
	config, err := config.NewWithPath(path, 0x1B)
	common.Assert(err)

	cac := cache.NewCache(config.Redis.REDIS_IP, config.Redis.REDIS_PORT, config.Redis.REDIS_PREFIX)
	model, err := database.NewDao(config.Database.DB_URL, config.Database.DB_USER, config.Database.DB_PW, config.Database.DB_NAME)
	common.Assert(err)

	model.SetMaxIdleConns(config.Database.MAX_IDLE_NUM)
	model.SetMaxOpenConns(config.Database.MAX_OPEN_NUM)

	vr := view.NewViewRender(config.Website.TEMP_FOLDER)
	vr.SetDefaultVar(config.Website.SITE_NAME, config.Server.SITE_URL, config.Website.SITE_DES, config.Website.SITE_FOOTER, config.Website.SITE_KEYWORDS)

	context := controller.NewContext()
	context.Cache = cac
	context.Config = config
	context.Model = model
	context.View = vr

	if isRelease {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Static("/static", config.Server.STATIC_FILE_FOLDER)
	router.StaticFile("/b.css", config.Server.STATIC_FILE_FOLDER+"/b.css")
	router.StaticFile("/favicon.ico", config.Server.STATIC_FILE_FOLDER+"/favicon.ico")

	router.GET(net.SERVER_INDEX, context.Ctrl.IndC.Do)
	router.GET(net.SERVER_TWEET, context.Ctrl.TweC.Do)
	router.GET(net.SERVER_PAGE, context.Ctrl.PagC.Do)
	router.GET(net.SERVER_CAT, context.Ctrl.CatC.Do)
	router.GET(net.SERVER_ART, context.Ctrl.ArtC.Do)

	router.NoMethod(context.Ctrl.ErrC.Do404)
	router.NoRoute(context.Ctrl.ErrC.Do404)

	router.Run(config.Server.URL + config.Server.PORT)
}
