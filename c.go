package beaker

type BaseMsg struct {
	Code int `json:"code"`
	Data interface{}
}

type SuccMsg struct {
	BaseMsg
	Msg string `json:"msg"`
}

type CtrlServer struct {
	mvc    *mvc
	ArcC   *ArcCtrl
	CatC   *CatCtrl
	PagC   *PagCtrl
	TweC   *TweCtrl
	ErrC   *ErrCtrl
	IndC   *IndCtrl
	OptC   *OptCtrl
	LoginC *LoginCtrl
}

type mvc struct {
	cache  *Cache
	model  *Model
	view   *ViewRender
	config *ConfigData
}

func NewCtrl() *CtrlServer {
	c := new(CtrlServer)
	c.ArcC = new(ArcCtrl)
	c.ArcC.ctrl = c
	c.CatC = new(CatCtrl)
	c.CatC.ctrl = c
	c.PagC = new(PagCtrl)
	c.PagC.ctrl = c
	c.TweC = new(TweCtrl)
	c.TweC.ctrl = c
	c.ErrC = new(ErrCtrl)
	c.ErrC.ctrl = c
	c.IndC = new(IndCtrl)
	c.IndC.ctrl = c
	c.OptC = new(OptCtrl)
	c.OptC.ctrl = c
	c.LoginC = new(LoginCtrl)
	c.LoginC.ctrl = c
	c.mvc = new(mvc)
	return c
}

func (c *CtrlServer) SetCache(cache *Cache) {
	c.mvc.cache = cache
}

func (c *CtrlServer) SetModel(model *Model) {
	c.mvc.model = model
}

func (c *CtrlServer) SetViewRender(vr *ViewRender) {
	c.mvc.view = vr
}

func (c *CtrlServer) SetConfig(conf *ConfigData) {
	c.mvc.config = conf
}

type BaseCtrl struct {
	ctrl *CtrlServer
}

type ArcCtrl struct {
	BaseCtrl
}

type CatCtrl struct {
	BaseCtrl
}

type PagCtrl struct {
	BaseCtrl
}

type TweCtrl struct {
	BaseCtrl
}

type ErrCtrl struct {
	BaseCtrl
}

type IndCtrl struct {
	BaseCtrl
}

type OptCtrl struct {
	BaseCtrl
}

type LoginCtrl struct {
	BaseCtrl
}
