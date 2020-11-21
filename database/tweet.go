package database

import "github.com/mebiusashan/beaker/common"

type TweetModelDB struct {
	common.TweetModel
}

func (TweetModelDB) TableName() string {
	return "tweets"
}

func (d *dao) TweetFindByNum(page uint, num uint) ([]TweetModelDB, error) {
	var tw []TweetModelDB
	err := d.mysql.Order("id desc").Offset((page - 1) * num).Limit(num).Find(&tw).Error
	return tw, err
}

func (d *dao) TweetCount() uint {
	var count uint = 0
	d.mysql.Table("tweets").Count(&count)
	return count
}

func (d *dao) TweetAdd(content string) error {
	tw := TweetModelDB{}
	tw.Content = content
	err := d.mysql.Create(&tw).Error
	return err
}

func (d *dao) TweetDel(id uint) error {
	tw := TweetModelDB{}
	tw.ID = id
	err := d.mysql.Delete(tw).Error
	return err
}
