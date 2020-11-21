package controller

import (
	"github.com/mebiusashan/beaker/cache"
	"github.com/mebiusashan/beaker/config"
	"github.com/mebiusashan/beaker/database"
	"github.com/mebiusashan/beaker/view"
)

type Context struct {
	Config *config.ConfigData
	Cache  *cache.Cache
	View   *view.ViewRender
	Model  *database.Dao
	Ctrl   *Controller
}

type Controller struct {
	ArtC   *ArticleController
	CatC   *CategoryController
	PagC   *PageController
	TweC   *TweetController
	ErrC   *ErrerController
	IndC   *IndexController
	OptC   *OptionController
	LoginC *LoginController
}

type BaseController struct {
	Context *Context
}

func NewContext() *Context {
	c := new(Context)
	c.Ctrl = new(Controller)
	c.Ctrl.ArtC = new(ArticleController)
	c.Ctrl.CatC = new(CategoryController)
	c.Ctrl.PagC = new(PageController)
	c.Ctrl.TweC = new(TweetController)
	c.Ctrl.ErrC = new(ErrerController)
	c.Ctrl.IndC = new(IndexController)
	c.Ctrl.OptC = new(OptionController)
	c.Ctrl.LoginC = new(LoginController)
	return c
}
