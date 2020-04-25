package handler

import (
	cfg "filestore-server/config"
	dblayer "filestore-server/db"
	"filestore-server/meta"
	"filestore-server/store/localS3"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// DownloadURLHandler : 生成文件的下载地址
func DownloadURLHandler(c *gin.Context) {
	filehash := c.Request.FormValue("filehash")
	username := c.Request.FormValue("username")
	token := c.Request.FormValue("token")
	tmpURL := fmt.Sprintf(
		"http://%s/file/download?filehash=%s&username=%s&token=%s",
		c.Request.Host, filehash, username, token)
	c.Data(http.StatusOK, "octet-stream",[]byte(tmpURL))
}

// DownloadHandler : 文件下载接口
func DownloadHandler(c *gin.Context) {
	fsha1 := c.Request.FormValue("filehash")
	username := c.Request.FormValue("username")

	fm, _ := meta.GetFileMetaDB(fsha1)
	userFile, _ := dblayer.QueryUserFileMeta(username, fsha1)

	if strings.HasPrefix(fm.FileLocation, cfg.TempLocalRootDir) {
		c.FileAttachment(fm.FileLocation, userFile.FileName)
	} else if strings.HasPrefix(fm.FileLocation, "/files") {
		fmt.Println("to download file from s3...")
		bucket := localS3.GetS3Bucket(cfg.S3Bucker)
		data, _ := bucket.Get(fm.FileLocation)
		c.Header("content-type", "application/octect-stream")
		// attachment表示文件将会提示下载到本地，而不是直接在浏览器中打开
		c.Header("content-disposition", "attachment; filename=\""+userFile.FileName+"\"")
		c.Data(http.StatusOK, "application/octect-stream", data)
	}
}

// RangeDownloadHandler : 支持断点的文件下载接口
func RangeDownloadHandler(c *gin.Context) {
	fsha1 := c.Request.FormValue("filehash")
	username := c.Request.FormValue("username")

	fm, _ := meta.GetFileMetaDB(fsha1)
	userFile, err := dblayer.QueryUserFileMeta(username, fsha1)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}

	f, err := os.Open(fm.FileLocation)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	data, _ := ioutil.ReadAll(f)
	defer f.Close()

	c.Header("Content-Type", "application/octect-stream")
	// attachment表示文件将会提示下载到本地，而不是直接在浏览器中打开
	c.Header("content-disposition", "attachment; filename=\""+userFile.FileName+"\"")
	c.Data(http.StatusOK, "application/octect-stream", data)
}
