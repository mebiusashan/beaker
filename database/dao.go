package database

import "github.com/jinzhu/gorm"

type Dao struct {
	mysql  *gorm.DB
	server string
}

func NewDao(url string, name string, password string, dbname string) (*Dao, error) {
	d := new(Dao)
	d.server = name + ":" + password + "@tcp(" + url + ")/" + dbname +
		"?charset=utf8mb4&parseTime=true&loc=Local"
	return d, d.Open()
}

func (d *Dao) Open() error {
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

func (d *Dao) SetMaxIdleConns(n int) {
	d.mysql.DB().SetMaxIdleConns(n)
}

func (d *Dao) SetMaxOpenConns(n int) {
	d.mysql.DB().SetMaxOpenConns(n)
}
