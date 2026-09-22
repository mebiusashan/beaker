package controller

import (
	"math"

	"github.com/gin-gonic/gin"
	"github.com/mebiusashan/beaker/internal/common"
)

func (ct *TweetController) Do(c *gin.Context) {
	numOfOnePage := ct.Context.Config.Website.TWEET_NUM_ONE_PAGE
	twCount, err := ct.Context.Model.TweetCount()
	if hasErrDo500(c, ct.Context.Ctrl.ErrC, err) {
		return
	}
	totalPages := uint(math.Ceil(float64(twCount) / float64(numOfOnePage)))
	twpages := ct.createPageNum(totalPages)
	page := ct.convParam(c.Param("page"))
	if totalPages == 0 {
		page = 1
	} else if page > totalPages {
		page = totalPages
	}

	tws, err := ct.Context.Model.TweetFindByNum(page, numOfOnePage)
	if hasErrDo500(c, ct.Context.Ctrl.ErrC, err) {
		return
	}

	pages, err := ct.Context.Model.PageFindAll()
	if hasErrDo404(c, ct.Context.Ctrl.ErrC, err) {
		return
	}

	cats, err := ct.Context.Model.CategoryFindAll()
	if hasErrDo404(c, ct.Context.Ctrl.ErrC, err) {
		return
	}

	vars := ct.Context.View.GetVarMap()
	vars.Set("title", "Home")
	vars.Set("cats", cats)
	vars.Set("pages", pages)
	vars.Set("tws", tws)
	vars.Set("twpages", twpages)
	vars.Set("twcurpage", page)

	bodyStr, err := ct.Context.View.Render(common.TEMPLATE_TWEET, vars)
	if hasErrDo404(c, ct.Context.Ctrl.ErrC, err) {
		return
	}
	write200(c, bodyStr)
}
