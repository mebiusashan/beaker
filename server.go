package beaker

import (
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/cert"
	"github.com/mebiusashan/beaker/common"
	"github.com/mebiusashan/beaker/net"
)

func RunServer() {
	path := os.Getenv(common.SERVER_ENV_KEY)
	config, err := NewWithPath(path, 0x1B)
	common.Assert(err)

	cac := NewCache(config.Redis.REDIS_IP, config.Redis.REDIS_PORT, config.Redis.REDIS_PREFIX)
	model, err := NewDatabase(config.Database.DB_URL, config.Database.DB_USER, config.Database.DB_PW, config.Database.DB_NAME)
	common.Assert(err)

	model.SetMaxIdleConns(config.Database.MAX_IDLE_NUM)
	model.SetMaxOpenConns(config.Database.MAX_OPEN_NUM)

	vr := NewViewRender(config.Website.TEMP_FOLDER)
	vr.SetDefaultVar(config.Website.SITE_NAME, config.Website.SITE_URL, config.Website.SITE_DES, config.Website.SITE_FOOTER, config.Website.SITE_KEYWORDS)

	ctr := NewCtrl()
	ctr.SetCache(cac)
	ctr.SetModel(model)
	ctr.SetViewRender(vr)
	ctr.SetConfig(config)

	router := gin.Default()

	//router.StaticFS("/static", http.Dir(config.Website.STATIC_FILE_FOLDER))
	router.StaticFile("/b.css", config.Website.STATIC_FILE_FOLDER+"/b.css")
	router.StaticFile("/favicon.ico", config.Website.STATIC_FILE_FOLDER+"/favicon.ico")

	router.GET(net.SERVER_INDEX, ctr.IndC.Do)
	router.GET(net.SERVER_TWEET, ctr.TweC.Do)
	router.GET(net.SERVER_PAGE, ctr.PagC.Do)
	router.GET(net.SERVER_CAT, ctr.CatC.Do)
	router.GET(net.SERVER_ART, ctr.ArcC.Do)

	router.NoMethod(ctr.ErrC.Do404)
	router.NoRoute(ctr.ErrC.Do404)

	router.Run(config.Server.URL + config.Server.PORT)
}

func RunAdmin() {
	path := os.Getenv(common.ADMIN_ENV_KEY)
	config, err := NewWithPath(path, 0x1D)
	common.Assert(err)

	cac := NewCache(config.Redis.REDIS_IP, config.Redis.REDIS_PORT, config.Redis.REDIS_PREFIX)
	model, err := NewDatabase(config.Database.DB_URL, config.Database.DB_USER, config.Database.DB_PW, config.Database.DB_NAME)
	common.Assert(err)
	model.SetMaxIdleConns(config.Database.MAX_IDLE_NUM)
	model.SetMaxOpenConns(config.Database.MAX_OPEN_NUM)

	ctr := NewCtrl()
	ctr.SetCache(cac)
	ctr.SetModel(model)
	ctr.SetConfig(config)

	pubHas, err := PathExists(config.AuthInfo.ServerKeyDir + SERVER_PUBLIC_KEY)
	common.Assert(err)

	priHas, err := PathExists(config.AuthInfo.ServerKeyDir + SERVER_PRIVATE_KEY)
	common.Assert(err)

	if !pubHas || !priHas {
		pub, pri, err := cert.CreateRSAKeys()
		common.Assert(err)
		err = ioutil.WriteFile(config.AuthInfo.ServerKeyDir+SERVER_PUBLIC_KEY, pub, 0666)
		common.Assert(err)
		err = ioutil.WriteFile(config.AuthInfo.ServerKeyDir+SERVER_PRIVATE_KEY, pri, 0666)
		common.Assert(err)
	}

	pub, err := ioutil.ReadFile(config.AuthInfo.ServerKeyDir + SERVER_PUBLIC_KEY)
	common.Assert(err)

	pri, err := ioutil.ReadFile(config.AuthInfo.ServerKeyDir + SERVER_PRIVATE_KEY)
	common.Assert(err)

	rel := cert.CheckRSAKey(pub, pri)
	if !rel {
		common.Err("Secret key verification failed")
	}

	router := gin.Default()

	//登录
	user := router.Group(net.ADMIN_GROUP_USER)
	user.POST(net.ADMIN_PING, ctr.LoginC.Ping)
	user.POST(net.ADMIN_LOGIN, ctr.LoginC.Login)
	user.POST(net.ADMIN_CHECK, ctr.LoginC.Check)

	adminr := router.Group(net.ADMIN_GROUP_ADMIN)
	adminr.POST(net.ADMIN_ART_ADD, ctr.ArcC.Add)
	adminr.POST(net.ADMIN_ART_RM, ctr.ArcC.Del)
	adminr.POST(net.ADMIN_ART_LIST, ctr.ArcC.All)
	adminr.POST(net.ADMIN_ART_DOWNLOAD, ctr.ArcC.Down)

	adminr.POST(net.ADMIN_PAGE_ADD, ctr.PagC.Add)
	adminr.POST(net.ADMIN_PAGE_RM, ctr.PagC.Del)
	adminr.POST(net.ADMIN_PAGE_LIST, ctr.PagC.List)
	adminr.POST(net.ADMIN_PAGE_DOWNLOAD, ctr.PagC.Down)

	adminr.POST(net.ADMIN_TWEET_ADD, ctr.TweC.Add)
	adminr.POST(net.ADMIN_TWEET_RM, ctr.TweC.Del)
	adminr.POST(net.ADMIN_TWEET_LIST, ctr.TweC.List)

	adminr.POST(net.ADMIN_CAT_ADD, ctr.CatC.Add)
	adminr.POST(net.ADMIN_CAT_RM, ctr.CatC.Del)
	adminr.POST(net.ADMIN_CAT_LIST, ctr.CatC.All)
	adminr.POST(net.ADMIN_CAT_MODIFY, ctr.CatC.Update)

	adminr.POST(net.ADMIN_OPT, ctr.OptC.Info)
	adminr.POST(net.ADMIN_CLEAN, ctr.OptC.ClearCache)

	router.Run(config.Server.PORT)
}
