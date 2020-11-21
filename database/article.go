package database

import "github.com/mebiusashan/beaker/common"

type ArticleModelDB struct {
	common.ArticleModel
}

func (ArticleModelDB) TableName() string {
	return "arts"
}

func (d *dao) ArticleFindWithNum(num uint) ([]ArticleModelDB, error) {
	var arts []ArticleModelDB
	err := d.mysql.Order("id desc").Select("id, title, created_at").Limit(num).Find(&arts).Error
	return arts, err
}

func (d *dao) ArticleFindByID(id uint) (*ArticleModelDB, error) {
	art := new(ArticleModelDB)
	err := d.mysql.Where("id = ?", id).First(art).Error
	return art, err
}

func (d *dao) ArticleFindByCatID(catid uint) ([]ArticleModelDB, error) {
	var arts []ArticleModelDB
	err := d.mysql.Order("id desc").Select("id, title, created_at").Where("catid = ?", catid).Find(&arts).Error
	return arts, err
}

func (d *dao) ArticleFindAll() ([]ArticleModelDB, error) {
	var arts []ArticleModelDB
	err := d.mysql.Select("id, title, created_at, updated_at, catid").Find(&arts).Error
	return arts, err
}

func (d *dao) ArticleAdd(catid uint, title string, content string) error {
	art := new(ArticleModelDB)
	art.Catid = catid
	art.Title = title
	art.Content = content
	err := d.mysql.Create(&art).Error
	return err
}

func (d *dao) ArticleDel(id uint) error {
	art := new(ArticleModelDB)
	art.ID = id
	err := d.mysql.Delete(&art).Error
	return err
}

func (d *dao) ArticleUpdateCat(catID uint, mvCatId uint) error {
	art := new(ArticleModelDB)
	err := d.mysql.Model(&art).Where("catid = ?", catID).Update("catid", mvCatId).Error
	return err
}
