package handler

import (
	"filestore-server/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"

	rPool "filestore-server/cache/redis"
	dblayer "filestore-server/db"
)

// MultipartUploadInfo : 初始化信息
type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
}

const (
	ChunkDir = "/Users/zhouliren/git/data/chunks/"
	MergeDir = "/Users/zhouliren/git/data/merge/"
)

func init() {
	if err := os.MkdirAll(ChunkDir, 0744); err != nil {
		fmt.Println("无法指定目录用于存储分块文件: " + ChunkDir)
		os.Exit(1)
	}

	if err := os.MkdirAll(MergeDir, 0744); err != nil {
		fmt.Println("无法指定目录用于存储合并后文件: " + MergeDir)
		os.Exit(1)
	}
}

// InitialMultipartUploadHandler : 初始化分块上传
func InitialMultipartUploadHandler(c *gin.Context) {
	// 1. 解析用户请求参数
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filesize, err := strconv.Atoi(c.Request.FormValue("filesize"))
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": -1,
				"msg":  "params invalid",
			})
		return
	}

	// 2. 获得redis的一个连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 3. 生成分块上传的初始化信息
	upInfo := MultipartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024, // 5MB
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))),
	}


	// 4. 将初始化信息写入到redis缓存
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "chunkcount", upInfo.ChunkCount)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filehash", upInfo.FileHash)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filesize", upInfo.FileSize)

	// 5. 将响应初始化数据返回到客户端
	c.JSON(
		http.StatusOK,
		gin.H{
			"code": 0,
			"msg":  "OK",
			"data": upInfo,
		})
}

// UploadPartHandler : 上传文件分块
func UploadPartHandler(c *gin.Context) {

	//	username := r.Form.Get("username")
	uploadID := c.Request.FormValue("uploadid")
	chunkIndex := c.Request.FormValue("index")

	// 2. 获得redis连接池中的一个连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 3. 获得文件句柄，用于存储分块内容
	fpath := ChunkDir + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": 0,
				"msg":  "Upload part failed",
				"data": nil,
			})
		return
	}
	defer fd.Close()

	// 读取带宽 5M
	buf := make([]byte, 5 * 1024*1024)
	for {
		n, err := c.Request.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	// 4. 更新redis缓存状态
	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)

	// 5. 返回处理结果到客户端
	c.JSON(
		http.StatusOK,
		gin.H{
			"code": 0,
			"msg":  "OK",
			"data": nil,
		})}

// CompleteUploadHandler : 通知上传合并
func CompleteUploadHandler(c *gin.Context) {
	// 1. 解析请求参数
	uploadID := c.Request.FormValue("uploadid")
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filesize := c.Request.FormValue("filesize")
	filename := c.Request.FormValue("filename")

	// 2. 获得redis连接池中的一个连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 3. 通过uploadid查询redis并判断是否所有分块上传完成
	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+uploadID))
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": -1,
				"msg":  "OK",
				"data": nil,
			})
		return
	}
	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount++
		}
	}
	if totalCount != chunkCount {
		fmt.Printf("chuncks count invalid: %d %d\n", totalCount, chunkCount)
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": -2,
				"msg":  "OK",
				"data": nil,
			})
		return
	}

	if mergeSuc := util.MergeChuncksByShell(ChunkDir+uploadID, MergeDir+filehash, filehash); !mergeSuc {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"code": -3,
				"msg": "complete upload failed",
				"data": nil,
			})
		//w.Write(util.NewRespMsg(-3, "complete upload failed", nil).JSONBytes())
		return
	}

	// 5. 更新唯一文件表及用户文件表
	fsize, _ := strconv.Atoi(filesize)
	dblayer.OnFileUploadFinished(filehash, filename, int64(fsize), MergeDir+filehash)
	dblayer.OnUserFileUploadFinished(username, filehash, filename, int64(fsize))

	delRes := util.RemovePathByShell(ChunkDir + uploadID)
	if !delRes {
		fmt.Printf("Failed to delete chuncks as upload comoleted, uploadID: %s\n", uploadID)
	}

	// 6. 响应处理结果
	c.JSON(
		http.StatusOK,
		gin.H{
			"code": 0,
			"msg":  "OK",
			"data": nil,
		})
}

// CancelUploadHandler : 通知取消上传
func CancelUploadHandler(c *gin.Context) {
	// 1. 解析用户请求参数
	username := c.Request.FormValue("username")
	uploadID := c.Request.FormValue("uploadid")
	if len(uploadID) <= len(username) {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "params invalid",
			Data: nil,
		}
		c.Data(http.StatusInternalServerError, "application/json", resp.JSONBytes())
		return
	}


	// 2. 获得redis的一个连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 3. 检查uploadID是否存在, 并尝试删除upload信息
	data, err := rConn.Do("del", "MP_"+uploadID)
	if err != nil {
		resp := util.RespMsg{
			Code: -2,
			Msg: "cancel upload failed",
			Data: nil,
		}
		c.Data(http.StatusInternalServerError, "application/json", resp.JSONBytes())
		return
	}

	// 4. 删除已上传的分块文件
	delRes := util.RemovePathByShell(ChunkDir + uploadID)
	if !delRes {
		fmt.Printf("Failed to delete chuncks as upload canceld, uploadID: %s\n", uploadID)
	}

	if res, ok := data.(int64); !ok || res != 1 {
		resp := util.RespMsg{
			Code: -3,
			Msg: "cancel upload failed",
			Data: nil,
		}
		c.Data(http.StatusInternalServerError, "application/json", resp.JSONBytes())
		return
	}


	c.JSON(
		http.StatusOK,
		gin.H{
			"code": 0,
			"msg":  "OK",
			"data": nil,
		})
}
