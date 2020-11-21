package database

import "github.com/mebiusashan/beaker/common"

type CategoryModelDB struct {
	common.CatModel
}

func (CategoryModelDB) TableName() string {
	return "cats"
}

func (d *dao) CategoryFindAll() ([]CategoryModelDB, error) {
	var cats []CategoryModelDB
	err := d.mysql.Order("id desc").Find(&cats).Error
	return cats, err
}

func (d *dao) CategoryFindByAlias(alias string) (*CategoryModelDB, error) {
	cat := new(CategoryModelDB)
	err := d.mysql.Order("id desc").Where("alias = ?", alias).First(cat).Error
	return cat, err
}

func (d *dao) CategoryAdd(name string, alias string) error {
	cat := CategoryModelDB{}
	cat.Name = name
	cat.Alias = alias
	err := d.mysql.Create(&cat).Error
	return err
}

func (d *dao) CategoryDel(id uint) error {
	cat := CategoryModelDB{}
	cat.ID = id
	err := d.mysql.Delete(cat).Error
	return err
}

func (d *dao) CategoryUpdate(id uint, name string, alias string) error {
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
