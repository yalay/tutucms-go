package controllers

import (
	"common"
	"conf"
	"log"
	"models"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/kataras/iris"
)

var sqliteDb *models.MyDb

func init() {
	sqliteDb = models.NewMyDb()
	err := sqliteDb.OpenDataBase("sqlite3", "tutu.db")
	if err != nil {
		log.Panicf("open db err:%v\n", err)
	}
}

func ArticlesGetHandler(ctx *iris.Context) {
	articleId := common.Atoi32(ctx.Param("id"))
	if articleId == 0 {
		ctx.NotFound()
		return
	}

	var article = models.Article{}
	sqliteDb.DB.First(&article, articleId)
	if article.Title == "" {
		ctx.NotFound()
		return
	}

	ctx.JSON(iris.StatusOK, article)
}

func ArticlesPostHandler(ctx *iris.Context) {
	title := ctx.FormValue("title")
	if title == "" {
		ctx.Writef("title is empty")
		return
	}

	var count int
	sqliteDb.DB.Where("title = ?", title).Find(&models.Article{}).Count(&count)
	if count > 0 {
		ctx.Writef("title:%s exist. Please delete it first.", title)
		return
	}

	article := models.Article{
		Star:    1,
		Status:  1,
		Author:  conf.GetAuthor(),
		Comeurl: conf.GetComeUrl(),
	}

	tags := ctx.FormValue("tags")
	article.Cid = common.Atoi32(ctx.FormValue("cid"))
	article.Title = title
	article.Tag = tags
	article.Color = ctx.FormValue("color")
	article.Cover = ctx.FormValue("cover")
	article.Remark = ctx.FormValue("remark")
	article.ShortTitle = ctx.FormValue("short_title")
	article.Keywords = func() string {
		keywords := ctx.FormValue("keywords")
		if keywords != "" {
			return keywords
		}
		return GetKeywords(title)
	}()
	article.Content = ctx.FormValue("content")
	article.Addtime = int32(time.Now().Unix())
	sqliteDb.DB.Save(&article)

	// 保存tag
	if tags != "" {
		associton := sqliteDb.DB.Model(&article).Association("Tags")
		if associton == nil || associton.Error != nil {
			ctx.NotFound()
			return
		}
		tagFields := strings.Split(tags, ",")
		modelTags := make([]models.Tag, 0, len(tagFields))
		for _, tag := range tagFields {
			if tag == "" {
				continue
			}
			modelTags = append(modelTags, models.Tag{
				Tag:   tag,
				Title: title,
			})
		}
		associton.Append(modelTags)
	}
	ctx.Writef("add success")
}

func AttachsGetHandler(ctx *iris.Context) {
	articleId := common.Atoi32(ctx.Param("id"))
	if articleId == 0 {
		ctx.NotFound()
		return
	}

	associton := sqliteDb.DB.Model(&models.Article{ID: articleId}).Association("Attachs")
	if associton == nil || associton.Error != nil {
		ctx.NotFound()
		return
	}

	attachs := make([]models.Attach, associton.Count())
	associton.Find(&attachs)
	if len(attachs) == 0 {
		ctx.NotFound()
		return
	}

	ctx.JSON(iris.StatusOK, attachs)
}

func AttachsPostHandler(ctx *iris.Context) {
	articleId := common.Atoi32(ctx.Param("id"))
	if articleId == 0 {
		ctx.NotFound()
		return
	}

	file := ctx.FormValue("file")
	if file == "" {
		ctx.Writef("file is empty")
		return
	}

	var count int
	sqliteDb.DB.Where("file = ?", file).Find(&models.Attach{}).Count(&count)
	if count > 0 {
		ctx.Writef("file:%s exist. Please delete it first.", file)
		return
	}

	associton := sqliteDb.DB.Model(&models.Article{ID: articleId}).Association("Attachs")
	if associton == nil || associton.Error != nil {
		ctx.NotFound()
		return
	}

	attach := models.Attach{
		ArticleId: articleId,
		Uid:       conf.GetUserId(),
		Name:      filepath.Base(file),
		Remark:    ctx.FormValue("remark"),
		Size:      common.Atoi32(ctx.FormValue("size")),
		File:      file,
		Ext:       filepath.Ext(file),
		Type: func() int32 {
			if strings.HasPrefix(file, "http://") ||
				strings.HasPrefix(file, "https://") {
				return 1
			}
			return 0
		}(),
		Status:     1,
		UploadTime: int32(time.Now().Unix()),
	}
	associton.Append(attach)
	ctx.Writef("add success")
}

func TagsGetHandler(ctx *iris.Context) {
	articleId := common.Atoi32(ctx.Param("id"))
	if articleId == 0 {
		ctx.NotFound()
		return
	}

	associton := sqliteDb.DB.Model(&models.Article{ID: articleId}).Association("Tags")
	if associton == nil || associton.Error != nil {
		ctx.NotFound()
		return
	}

	tags := make([]models.Tag, associton.Count())
	associton.Find(&tags)
	if len(tags) == 0 {
		ctx.NotFound()
		return
	}

	ctx.JSON(iris.StatusOK, tags)
}
