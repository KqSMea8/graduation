package model

import "time"

type File struct {
	Uid        int64     `json:"uid";gorm:"uid"`
	Fid        int64     `json:"fid";gorm:"fid"`
	Name       string    `json:"name";gorm:"name"`
	Size       int64     `json:"size";gorm:"fid"`
	Md5        string    `json:"md5";gorm:"md5"`
	CreateTime time.Time `json:"create_time";gorm:"create_time"`
	UpdateTime time.Time `json:"update_time";gorm:"update_time"`
}

func (File) TableName() string {
	return "file_info"
}
