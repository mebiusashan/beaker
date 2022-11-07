package database

import "github.com/mebiusashan/beaker/internal/common"

type TweetModelDB struct {
	common.TweetModel
}

func (TweetModelDB) TableName() string {
	return "tweets"
}

func (d *Dao) TweetFindByNum(page uint, num uint) ([]TweetModelDB, error) {
	var tw []TweetModelDB
	count := 0
	d.mysql.Model(&TweetModelDB{}).Count(&count)
	if count == 0 {
		return tw, nil
	}
	err := d.mysql.Order("id desc").Offset((page - 1) * num).Limit(num).Find(&tw).Error
	return tw, err
}

func (d *Dao) TweetCount() uint {
	var count uint = 0
	d.mysql.Table("tweets").Count(&count)
	return count
}

func (d *Dao) TweetAdd(content string) error {
	tw := TweetModelDB{}
	tw.Content = content
	err := d.mysql.Create(&tw).Error
	return err
}

func (d *Dao) TweetDel(id uint) error {
	tw := TweetModelDB{}
	tw.ID = id
	err := d.mysql.Delete(tw).Error
	return err
}
