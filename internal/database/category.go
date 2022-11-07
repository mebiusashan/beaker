package database

import "github.com/mebiusashan/beaker/internal/common"

type CategoryModelDB struct {
	common.CatModel
}

func (CategoryModelDB) TableName() string {
	return "cats"
}

func (d *Dao) CategoryFindAll() ([]CategoryModelDB, error) {
	var cats []CategoryModelDB
	err := d.mysql.Order("id desc").Find(&cats).Error
	return cats, err
}

func (d *Dao) CategoryFindByAlias(alias string) (*CategoryModelDB, error) {
	cat := new(CategoryModelDB)
	err := d.mysql.Order("id desc").Where("alias = ?", alias).First(cat).Error
	return cat, err
}

func (d *Dao) CategoryFindByID(id uint) (*CategoryModelDB, error) {
	cat := new(CategoryModelDB)
	err := d.mysql.Order("id desc").Where("id = ?", id).First(cat).Error
	return cat, err
}

func (d *Dao) CategoryAdd(name string, alias string) error {
	cat := CategoryModelDB{}
	cat.Name = name
	cat.Alias = alias
	err := d.mysql.Create(&cat).Error
	return err
}

func (d *Dao) CategoryDel(id uint) error {
	cat := CategoryModelDB{}
	cat.ID = id
	err := d.mysql.Delete(cat).Error
	return err
}

func (d *Dao) CategoryUpdate(id uint, name string, alias string) error {
	data := make(map[string]interface{})
	if name != "" {
		data["name"] = name
	}
	if alias != "" {
		data["alias"] = alias
	}
	cat := CategoryModelDB{}
	cat.ID = id
	err := d.mysql.Model(&cat).Updates(data).Error
	return err
}
