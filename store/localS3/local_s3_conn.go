package localS3

import (
	cfg "filestore-server/config"
	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v1/s3"
)

var locals3 *s3.S3

func GetS3Connection() *s3.S3  {
	//init link info
	if(locals3 !=nil) {
		return locals3
	}

	auth := aws.Auth{
		AccessKey: cfg.S3AccessKey,
		SecretKey: cfg.S3SecretKey,
	}
	curRegion := aws.Region{
		Name:"default",
		EC2Endpoint: cfg.S3GEndpoint,
		S3Endpoint:cfg.S3GEndpoint,
		S3BucketEndpoint:"",
		S3LocationConstraint: false,
		//S3LowercaseBucket:false,
		Sign: aws.SignV2,
	}
	return s3.New(auth, curRegion)
}

// GetCephBucket : 获取指定的bucket对象
func GetS3Bucket(bucket string) *s3.Bucket {
	conn := GetS3Connection()
	return conn.Bucket(bucket)
}

func PutObject(bucket string, path string, data []byte) error {
	return GetS3Bucket(bucket).Put(path, data, "octet-stream",s3.PublicRead)
}