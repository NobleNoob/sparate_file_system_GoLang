package main

import (
	"filestore-server/store/localS3"
	"fmt"
	"gopkg.in/amz.v1/s3"
	//"os"
)

func main() {
	bucket := localS3.GetS3Bucket("1")

	res, _ := bucket.List("", "", "", 100)
	fmt.Printf("object keys: %+v\n", res)

	// 新上传一个对象
	err := bucket.Put("/users/zhouliren/git/a", []byte("just for test"), "octet-stream", s3.PublicRead)
	fmt.Printf("upload err: %+v\n", err)

	// 查询这个bucket下面指定条件的object keys
	res, err = bucket.List("", "", "", 100)
	fmt.Printf("object keys: %+v\n", res)
}