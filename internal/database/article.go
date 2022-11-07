package database

import (
	"github.com/mebiusashan/beaker/internal/common"
)

type ArticleModelDB struct {
	common.ArticleModel
}

func (ArticleModelDB) TableName() string {
	return "arts"
}

func (d *Dao) ArticleFindWithNum(num uint) ([]ArticleModelDB, error) {
	var arts []ArticleModelDB
	err := d.mysql.Order("id desc").Select("id, title, created_at").Limit(num).Find(&arts).Error
	return arts, err
}

func (d *Dao) ArticleFindByID(id uint) (*ArticleModelDB, error) {
	art := new(ArticleModelDB)
	err := d.mysql.Where("id = ?", id).First(art).Error
	return art, err
}

func (d *Dao) ArticleFindByCatID(catid uint) ([]ArticleModelDB, error) {
	var arts []ArticleModelDB
	err := d.mysql.Order("id desc").Select("id, title, created_at").Where("catid = ?", catid).Find(&arts).Error
	return arts, err
}

func (d *Dao) ArticleFindAll() ([]ArticleModelDB, error) {
	var arts []ArticleModelDB
	err := d.mysql.Select("id, title, created_at, updated_at, catid").Find(&arts).Error
	return arts, err
}

func (d *Dao) ArticleAdd(catid uint, title string, content string) error {
	art := new(ArticleModelDB)
	art.Catid = catid
	art.Title = title
	art.Content = content
	err := d.mysql.Create(&art).Error
	return err
}

func (d *Dao) ArticleDel(id uint) error {
	art := new(ArticleModelDB)
	art.ID = id
	return d.mysql.Delete(&art).Error
}

func (d *Dao) ArticleUpdateCat(catID uint, mvCatId uint) error {
	art := new(ArticleModelDB)
	return d.mysql.Model(&art).Where("catid = ?", catID).Update("catid", mvCatId).Error
}

func (d *Dao) ArticleUpdate(id uint, m *common.ArticleModel) error {
	art, err := d.ArticleFindByID(id)
	if err != nil {
		return err
	}
	if m.Catid != 0 {
		art.Catid = m.Catid
	}
	if m.Content != "" {
		art.Content = m.Content
	}
	if m.Title != "" {
		art.Title = m.Title
	}
	return d.mysql.Model(art).Updates(art).Error
}
