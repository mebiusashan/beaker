package database

import "github.com/mebiusashan/beaker/internal/common"

type PageModelDB struct {
	common.PageModel
}

func (PageModelDB) TableName() string {
	return "pages"
}

func (d *Dao) PageFindAll() ([]PageModelDB, error) {
	var pages []PageModelDB
	err := d.mysql.Order("id").Select("id, title").Find(&pages).Error
	return pages, err
}

func (d *Dao) PageFindByID(id uint) (*PageModelDB, error) {
	page := new(PageModelDB)
	err := d.mysql.Where("id = ?", id).First(page).Error
	return page, err
}

func (d *Dao) PageAdd(title string, content string) error {
	page := new(PageModelDB)
	page.Title = title
	page.Content = content
	err := d.mysql.Create(&page).Error
	return err
}

func (d *Dao) PageDel(id uint) error {
	page := PageModelDB{}
	page.ID = id
	err := d.mysql.Delete(&page).Error
	return err
}

func (d *Dao) PageUpdate(id uint, m *common.PageModel) error {
	page, err := d.PageFindByID(id)
	if err != nil {
		return err
	}
	if m.Content != "" {
		page.Content = m.Content
	}
	if m.Title != "" {
		page.Title = m.Title
	}
	return d.mysql.Model(page).Updates(page).Error
}
