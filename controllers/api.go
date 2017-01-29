package controllers

import (
	"common"
	"conf"
	"fmt"
	"log"
	"models"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kataras/iris"
)

var sqliteDb *models.MyDb

func init() {
	sqliteDb = models.NewMyDb()
	err := sqliteDb.OpenDataBase("mysql", "root:ya@/dbname?charset=utf8&parseTime=True&loc=Local")
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
	sqliteDb.DB.Model(&models.Article{}).Where("title = ?", title).Count(&count)
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

	titleKeywords := GetKeywords(title)
	article.Cid = common.Atoi32(ctx.FormValue("cid"))
	article.Title = title
	article.Tag = func() string {
		tags := ctx.FormValue("tags")
		if tags != "" {
			return tags
		}
		return titleKeywords
	}()
	article.Cover = ctx.FormValue("cover")
	article.Remark = title
	article.Keywords = func() string {
		keywords := ctx.FormValue("keywords")
		if keywords != "" {
			return keywords
		}
		return GetKeywords(title)
	}()
	article.Content = func() string {
		content := ctx.FormValue("content")
		if content != "" {
			return content
		}
		return title
	}()
	article.Addtime = int32(time.Now().Unix())

	// 保存tag
	if article.Tag != "" {
		tagFields := strings.Split(article.Tag, ",")
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
		article.Tags = append(article.Tags, modelTags...)
	}
	sqliteDb.DB.Save(&article)
	ctx.JSON(iris.StatusOK, article)
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
	files := ctx.FormValue("files")
	if file == "" && files == "" {
		ctx.Writef("file is empty")
		return
	}

	size := common.Atoi64(ctx.FormValue("size"))
	if file != "" {
		err := addAttachFile(file, articleId, size)
		if err != nil {
			ctx.Writef("%v", err)
			return
		}
	} else if files != "" {
		result := ""
		basePaths := filepath.Base(files)
		dirPath := filepath.Dir(files)
		fileFields := strings.Split(basePaths, ",")
		for _, fileField := range fileFields {
			err := addAttachFile(filepath.Join(dirPath, fileField), articleId, size)
			if err != nil {
				result += err.Error()
			}
		}
		if result != "" {
			ctx.Writef(result)
			return
		}
	}

	ctx.Writef("add success")
}

func addAttachFile(file string, articleId int32, size int64) error {
	var count int
	sqliteDb.DB.Where("file = ?", file).Find(&models.Attach{}).Count(&count)
	if count > 0 {
		return fmt.Errorf("file:%s exist. Please delete it first.", file)
	}

	var article = models.Article{}
	sqliteDb.DB.First(&article, articleId)
	if article.Title == "" {
		return fmt.Errorf("title is empty. file:%s.", file)
	}

	attach := models.Attach{
		ArticleId: articleId,
		Uid:       conf.GetUserId(),
		Name:      filepath.Base(file),
		Remark:    article.Title,
		Size: func() int64 {
			if size != 0 {
				return size
			}
			return getSize(file)
		}(),
		File: file,
		Ext:  filepath.Ext(file),
		Type: func() int32 {
			if isRemoteFile(file) {
				return 1
			}
			return 0
		}(),
		Status:     1,
		UploadTime: int32(time.Now().Unix()),
	}
	sqliteDb.DB.Save(&attach)
	return nil
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

func TagsPostHandler(ctx *iris.Context) {
	articleId := common.Atoi32(ctx.Param("id"))
	if articleId == 0 {
		ctx.NotFound()
		return
	}

	tags := ctx.FormValue("tags")
	if tags == "" {
		ctx.Writef("tags is empty")
		return
	}

	article := models.Article{ID: articleId}
	sqliteDb.DB.First(&article)
	if article.Title == "" {
		ctx.Writef("article:%d donot exist.", articleId)
		return
	}

	newTagFields := strings.Split(tags, ",")

	// 更新article的tags字段
	if article.Tag != "" {
		articleTagFields := strings.Split(article.Tag, ",")
		tagSet := common.NewSet()
		for _, articleTag := range articleTagFields {
			tagSet.Add(articleTag)
		}
		for _, newTag := range newTagFields {
			if tagSet.Contains(newTag) {
				continue
			}
			article.Tag = article.Tag + "," + newTag
		}
	}

	// 更新tag表
	result := ""
	for _, tag := range newTagFields {
		var count int
		sqliteDb.DB.Where("article_id = ? AND tag = ?", articleId, tag).Find(&models.Tag{}).Count(&count)
		if count > 0 {
			if result == "" {
				result = tag
			} else {
				result += "," + tag
			}

			continue
		}

		var modelsTag = models.Tag{
			Tag:       tag,
			ArticleId: articleId,
			Title:     article.Title,
		}
		sqliteDb.DB.Save(&modelsTag)
	}

	if result != "" {
		ctx.Writef(result + " conflict")
	} else {
		ctx.Writef("add success")
	}
}

func isRemoteFile(filePath string) bool {
	if strings.HasPrefix(filePath, "http://") ||
		strings.HasPrefix(filePath, "https://") {
		return true
	}
	return false
}

func getSize(filePath string) int64 {
	// 远程文件暂不支持
	if isRemoteFile(filePath) {
		return 0
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0
	}
	return fileInfo.Size()
}
