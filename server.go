package beaker

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/cert"
)

func RunServer() {
	path := os.Getenv("BEAKERPATH")
	config, err := NewWithPath(path, 0x1B)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cac := NewCache(config.Redis.REDIS_IP, config.Redis.REDIS_PORT, config.Redis.REDIS_PREFIX)

	model, err := NewDatabase(config.Database.DB_URL, config.Database.DB_USER, config.Database.DB_PW, config.Database.DB_NAME)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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
	router.StaticFile("/favicon.ico", config.Website.STATIC_FILE_FOLDER+"favicon.ico")

	router.GET("/", ctr.IndC.Do)
	router.GET("/tweet/*page", ctr.TweC.Do)
	router.GET("/page/:id", ctr.PagC.Do)
	router.GET("/category/:name", ctr.CatC.Do)
	router.GET("/archives/:id", ctr.ArcC.Do)

	router.NoMethod(ctr.ErrC.Do404)
	router.NoRoute(ctr.ErrC.Do404)

	router.Run(config.Server.URL + config.Server.PORT)
}

func RunAdmin() {

	path := os.Getenv("HBADMINPATH")
	config, err := NewWithPath(path, 0x1D)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cac := NewCache(config.Redis.REDIS_IP, config.Redis.REDIS_PORT, config.Redis.REDIS_PREFIX)
	model, err := NewDatabase(config.Database.DB_URL, config.Database.DB_USER, config.Database.DB_PW, config.Database.DB_NAME)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	model.SetMaxIdleConns(config.Database.MAX_IDLE_NUM)
	model.SetMaxOpenConns(config.Database.MAX_OPEN_NUM)

	ctr := NewCtrl()
	ctr.SetCache(cac)
	ctr.SetModel(model)
	ctr.SetConfig(config)

	pubHas, err := PathExists(config.AuthInfo.ServerKeyDir + SERVER_PUBLIC_KEY)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	priHas, err := PathExists(config.AuthInfo.ServerKeyDir + SERVER_PRIVATE_KEY)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(config.AuthInfo.ServerKeyDir, pubHas, priHas)

	if !pubHas || !priHas {
		pub, pri, err := cert.CreateRSAKeys()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = ioutil.WriteFile(config.AuthInfo.ServerKeyDir+SERVER_PUBLIC_KEY, pub, 0666)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = ioutil.WriteFile(config.AuthInfo.ServerKeyDir+SERVER_PRIVATE_KEY, pri, 0666)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	pub, err := ioutil.ReadFile(config.AuthInfo.ServerKeyDir + SERVER_PUBLIC_KEY)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pri, err := ioutil.ReadFile(config.AuthInfo.ServerKeyDir + SERVER_PRIVATE_KEY)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	succ := cert.CheckRSAKey(pub, pri)
	if !succ {
		fmt.Println(err)
		os.Exit(1)
	}

	router := gin.Default()

	//登录
	user := router.Group("user")
	user.POST("ping", ctr.LoginC.Ping)
	user.POST("login", ctr.LoginC.Login)
	user.POST("check", ctr.LoginC.Check)

	adminr := router.Group("admin")
	adminr.POST("arc/add", ctr.ArcC.Add)
	adminr.POST("arc/del", ctr.ArcC.Del)
	adminr.POST("arc/list", ctr.ArcC.All)
	adminr.POST("arc/down", ctr.ArcC.Down)

	adminr.POST("pag/add", ctr.PagC.Add)
	adminr.POST("pag/del", ctr.PagC.Del)
	adminr.POST("pag/list", ctr.PagC.List)
	adminr.POST("pag/down", ctr.PagC.Down)

	adminr.POST("twe/add", ctr.TweC.Add)
	adminr.POST("twe/del", ctr.TweC.Del)
	adminr.POST("twe/list", ctr.TweC.List)

	adminr.POST("cat/add", ctr.CatC.Add)
	adminr.POST("cat/del", ctr.CatC.Del)
	adminr.POST("cat/list", ctr.CatC.All)
	adminr.POST("cat/change", ctr.CatC.Update)

	adminr.POST("opt", ctr.OptC.Info)
	adminr.POST("clr/cache", ctr.OptC.ClearCache)

	router.Run(config.Server.PORT)
}
