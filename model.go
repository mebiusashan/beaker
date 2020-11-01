package beaker

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Model struct {
	blogDao *dao
	ArcDao  *ArcMDao
	CatDao  *CatMDao
	PagDao  *PagMDao
	TweDao  *TweMDao
}

type dao struct {
	mysql  *gorm.DB
	server string
}

type baseModel struct {
	blogDao *dao
}

//创建一个新的数据库操作器
func NewDatabase(url string, name string, password string, dbname string) (*Model, error) {
	m := new(Model)
	m.blogDao = new(dao)
	err := m.blogDao.newDao(url, name, password, dbname)
	if err != nil {
		return nil, err
	}
	m.ArcDao = newArc(m.blogDao)
	m.CatDao = newCat(m.blogDao)
	m.PagDao = newPage(m.blogDao)
	m.TweDao = newTweet(m.blogDao)

	return m, nil
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

func (m *Model) SetMaxIdleConns(n int) {
	m.blogDao.mysql.DB().SetMaxIdleConns(n)
}

func (m *Model) SetMaxOpenConns(n int) {
	m.blogDao.mysql.DB().SetMaxOpenConns(n)
}
