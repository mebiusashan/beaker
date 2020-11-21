package database

import "github.com/jinzhu/gorm"

type dao struct {
	mysql  *gorm.DB
	server string
}

func (d *dao) newDao(url string, name string, password string, dbname string) error {
	d.server = name + ":" + password + "@tcp(" + url + ")/" + dbname +
		"?charset=utf8mb4&parseTime=true&loc=Local"
	return d.open()
}

func (d *dao) open() error {
	var err error
	d.mysql, err = gorm.Open("mysql", d.server)
	if err != nil {
		return err
	}
	err = d.mysql.DB().Ping()
	if err != nil {
		return err
	}

	d.mysql.SingularTable(true)
	return nil
}
