package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	cmn "filestore-server/common"
	cfg "filestore-server/config"
	dblayer "filestore-server/db"
	"filestore-server/meta"
	"filestore-server/mq"
	"filestore-server/store/localS3"
	"filestore-server/util"
)

func init() {
	// 目录已存在
	if _, err := os.Stat(cfg.TempLocalRootDir); err == nil {
		return
	}

	// 尝试创建目录
	err := os.MkdirAll(cfg.TempLocalRootDir, 0744)
	if err != nil {
		log.Println("Tmp folder creation failed")
		os.Exit(1)
	}
}

func GetUploadHandler(c *gin.Context) {
	data, err := ioutil.ReadFile("./static/view/upload.html")
	if err != nil {
		c.String(404, `网页不存在`)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", data)
}


// UploadHandler ： 处理文件上传
func PostUploadHandler(c *gin.Context) {
		errCode := 0
		defer func() {
			if errCode < 0 {
				c.JSON(http.StatusOK, gin.H{
					"code": errCode,
					"msg":  "Upload failed",
				})
			}
		}()

		// 从Form 表单获取文件信息
		file, head, err := c.Request.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get form data, err:%s\n", err.Error())
			errCode = -1
			return
		}
		defer file.Close()

		// 把文件内容转为[]byte
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			fmt.Printf("Failed to get file data, err:%s\n", err.Error())
			errCode = -2
			return
		}

		// Init FileMeta
		fileMeta := meta.Filemeta{
			FileName: head.Filename,
			FileSha1: util.Sha1(buf.Bytes()),
			FileSize: int64(len(buf.Bytes())),
			UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
		}

		// Write filemeta to local tmp file
		fileMeta.FileLocation = cfg.TempLocalRootDir + fileMeta.FileSha1
		newFile, err := os.Create(fileMeta.FileLocation)
		if err != nil {
			fmt.Printf("Failed to create file, err:%s\n", err.Error())
			errCode = -3 // define error code for gin
			return
		}
		defer newFile.Close()

		nByte, err := newFile.Write(buf.Bytes())
		if int64(nByte) != fileMeta.FileSize || err != nil {
			fmt.Printf("Failed to save data into file, err:%s\n, writtenSize:%s", err.Error(),nByte)
			errCode = -4
			return
		}

		newFile.Seek(0, 0) //从头读取文件
		fileMeta.FileSha1 = util.FileSha1(newFile)

		newFile.Seek(0, 0)
		if cfg.CurrentStoreType == cmn.StoreS3 {
			data, _ := ioutil.ReadAll(newFile)
			s3Path := "/files/" + fileMeta.FileSha1
			if !cfg.AsyncTransferEnable {
				// 设置S3 bucket
				err = localS3.PutObject("test1", s3Path, data)
				if err != nil {
					fmt.Println(err.Error())
					errCode = -5
					return
				}
				fileMeta.FileLocation = s3Path
			} else {

				data := mq.TransferData{
					FileHash:      fileMeta.FileSha1,
					CurLocation:   fileMeta.FileLocation,
					DestLocation:  s3Path,
					DestStoreType: cmn.StoreS3,
				}
				pubData, _ := json.Marshal(data)
				pubSuc := mq.Publish(
					cfg.TransExchangeName,
					cfg.TransS3RoutingKey,
					pubData,
				)
				if !pubSuc {
					// TODO: 当前发送转移信息失败，稍后重试
				}
			}
		}
		// 更新文件记录至Mysql
		_ = meta.SetFileMetaDB(fileMeta)

		// 更新用户文件列表
		username := c.Request.FormValue("username")
		suc := dblayer.OnUserFileUploadFinished(username, fileMeta.FileSha1,
			fileMeta.FileName, fileMeta.FileSize)
		if suc {
			c.Redirect(http.StatusFound, "/static/view/home.html")
		} else {
			errCode = -6
		}
	}


// UploadSucHandler : 上传已完成
func UploadSucHandler(c *gin.Context) {
		c.JSON(http.StatusOK,
		gin.H{
			"code": 0,
			"msg":  "Upload Finish!",
		})
}

// GetFileMetaHandler : 获取文件元信息
func GetFileMetaHandler(c *gin.Context) {

	filehash := c.Request.FormValue("filehash")
	fMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"code": -2,
				"msg":  "Upload failed!",
			})
		return
	}

	if fMeta != nil {
		data, err := json.Marshal(fMeta)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{
					"code": -3,
					"msg":  "Upload failed!",
				})
			return
		}
		c.Data(http.StatusOK, "application/json", data)
	} else {
		c.JSON(http.StatusOK,
			gin.H{
				"code": -4,
				"msg":  "No such file",
			})
	}
}




// FileQueryHandler : 查询批量的文件元信息
func FileQueryHandler(c *gin.Context) {
	r.ParseForm()

	limitCnt, _ := strconv.Atoi(c.Request.FormValue("limit"))
	username := c.Request.FormValue("username")
	userFiles, err := dblayer.QueryUserFileMetas(username, limitCnt)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"code": -1,
				"msg":  "Query failed!",
			})
		return
	}

	data, err := json.Marshal(userFiles)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"code": -2,
				"msg":  "Query failed!",
			})
		return
	}
	c.Data(http.StatusOK, "application/json", data)
}

// FileMetaUpdateHandler ： 更新元信息接口(重命名)
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	username := r.Form.Get("username")
	newFileName := r.Form.Get("filename")

	if opType != "0" || len(newFileName) < 1 {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 更新用户文件表tbl_user_file中的文件名，tbl_file的文件名不用修改
	_ = dblayer.RenameFileName(username, fileSha1, newFileName)

	// 返回最新的文件信息
	userFile, err := dblayer.QueryUserFileMeta(username, fileSha1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(userFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// FileDeleteHandler : 删除文件及元信息
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	fileSha1 := r.Form.Get("filehash")

	// 删除本地文件
	fm, err := meta.GetFileMetaDB(fileSha1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	os.Remove(fm.FileLocation)

	// 删除文件表中的一条记录
	suc := dblayer.DeleteUserFile(username, fileSha1)
	if !suc {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// TryFastUploadHandler : 尝试秒传接口
func TryFastUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// 1. 解析请求参数
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filename := r.Form.Get("filename")
	filesize, _ := strconv.Atoi(r.Form.Get("filesize"))

	// 2. 从文件表中查询相同hash的文件记录
	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 3. 查不到记录则返回秒传失败
	if fileMeta.FileSha1 == "" {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "秒传失败，请访问普通上传接口",
		}
		w.Write(resp.JSONBytes())
		return
	}

	// 4. 上传过则将文件信息写入用户文件表， 返回成功
	suc := dblayer.OnUserFileUploadFinished(
		username, filehash, filename, int64(filesize))
	if suc {
		resp := util.RespMsg{
			Code: 0,
			Msg:  "秒传成功",
		}
		w.Write(resp.JSONBytes())
		return
	}
	resp := util.RespMsg{
		Code: -2,
		Msg:  "秒传失败，请稍后重试",
	}
	w.Write(resp.JSONBytes())
	return
}