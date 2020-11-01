package beaker

import (
	"bytes"
	"reflect"
	"time"

	"github.com/CloudyKit/jet"
)

type ViewRender struct {
	templateFolder string
	_SITE_NAME     string
	_SITE_URL      string
	_SITE_DES      string
	_SITE_FOOTER   string
}

func NewViewRender(templateFolder string) *ViewRender {
	vr := new(ViewRender)
	vr.templateFolder = templateFolder
	return vr
}

func (r *ViewRender) SetDefaultVar(siteName string, siteURL string, siteDes string, siteFooter string) {
	r._SITE_NAME = siteName
	r._SITE_URL = siteURL
	r._SITE_DES = siteDes
	r._SITE_FOOTER = siteFooter
}

func (r *ViewRender) GetVarMap() jet.VarMap {
	vars := make(jet.VarMap)
	vars.Set("SITE_NAME", r._SITE_NAME)
	vars.Set("SITE_URL", r._SITE_URL)
	vars.Set("SITE_DES", r._SITE_DES)
	vars.Set("SITE_FOOTER", r._SITE_FOOTER)
	vars.Set("ISINDEX", false)
	return vars
}

func (r *ViewRender) Render(templateName string, vars jet.VarMap) (string, error) {
	viewManage := jet.NewHTMLSet(r.templateFolder)
	t, err := viewManage.GetTemplate(templateName)
	if err != nil {
		// template could not be loaded
		return "", err
	}
	r.addTimepFunc(viewManage)
	r.addTTimepFunc(viewManage)
	var w bytes.Buffer
	if err = t.Execute(&w, vars, nil); err != nil {
		// error when executing template
		return "", err
	}

	return w.String(), nil
}

func (r *ViewRender) addTimepFunc(viewManage *jet.Set) {
	viewManage.AddGlobalFunc("timep", func(arguments jet.Arguments) reflect.Value {
		var t time.Time
		t = arguments.Get(0).Interface().(time.Time)
		str := t.Format("2006.01.02")
		//fmt.Println("时间"+str+"结束")
		return reflect.ValueOf(str)
	})
}

func (r *ViewRender) addTTimepFunc(viewManage *jet.Set) {
	viewManage.AddGlobalFunc("ttimep", func(arguments jet.Arguments) reflect.Value {
		var t time.Time
		t = arguments.Get(0).Interface().(time.Time)
		str := t.Format("15:04:05 2006-01-02")
		//fmt.Println("时间"+str+"结束")
		return reflect.ValueOf(str)
	})
}
