package models

import (
	"github.com/jinzhu/gorm"
)

const (
	kTablePrefix = "tutu_"
)

type MyDb struct {
	*gorm.DB
}

func NewMyDb() *MyDb {
	return &MyDb{}
}

type Article struct {
	ID         int32  `gorm:"primary_key"`
	Cid        int32  `gorm:"not null"`
	Title      string `gorm:"not null;index"`
	Tag        string
	Color      string   `gorm:"type:varchar(8)"`
	Cover      string   `gorm:"type:varchar(250)"`
	Author     string   `gorm:"type:varchar(50)"`
	Comeurl    string   `gorm:"type:varchar(250)"`
	Remark     string   `gorm:"type:text;not null"`
	ShortTitle string   `gorm:"type:text"`
	Keywords   string   `gorm:"type:text"`
	Content    string   `gorm:"type:longtext"`
	Hits       int32    `gorm:"not null"`
	Star       int32    `gorm:"not null"` // 默认1
	Status     int32    `gorm:"not null"` // 默认1
	Up         int32    `gorm:"not null"`
	Down       int32    `gorm:"not null"`
	Addtime    int32    `gorm:"not null"`
	Attachs    []Attach `gorm:"ForeignKey:ArticleId;AssociationForeignKey:ID"`
	Tags       []Tag    `gorm:"ForeignKey:ArticleId;AssociationForeignKey:ID"`
}

func (Article) TableName() string {
	return kTablePrefix + "article"
}

type Attach struct {
	ID         int32 `gorm:"primary_key"`
	ArticleId  int32
	Uid        int32
	Name       string `gorm:"type:varchar(100)"`
	Remark     string `gorm:"type:text;not null"`
	Size       int64
	File       string `gorm:"type:varchar(250);index"`
	Ext        string `gorm:"type:varchar(10)"`
	Status     int32  `gorm:"not null"` // 默认1 状态, 1:正常 0:隐藏
	Type       int32  `gorm:"not null"` // 附件类型, 0:本地文件, 1:网络文件
	TryCount   int32
	UploadTime int32 `gorm:"not null"`
}

func (Attach) TableName() string {
	return kTablePrefix + "attach"
}

type Tag struct {
	Tag       string `gorm:"type:varchar(30)"`
	Title     string `gorm:"type:varchar(250)"`
	ArticleId int32
}

func (Tag) TableName() string {
	return kTablePrefix + "tags"
}

func (m *MyDb) OpenDataBase(dbType, dbInfo string) error {
	myDb, err := gorm.Open(dbType, dbInfo)
	if err != nil {
		return err
	}

	if !myDb.HasTable(&Article{}) {
		myDb.CreateTable(&Article{})
	}
	if !myDb.HasTable(&Attach{}) {
		myDb.CreateTable(&Attach{})
	}
	if !myDb.HasTable(&Tag{}) {
		myDb.CreateTable(&Tag{})
	}
	myDb.Model(&Article{}).Related(&Attach{})
	myDb.Model(&Article{}).Related(&Tag{})
	m.DB = myDb
	return nil
}
