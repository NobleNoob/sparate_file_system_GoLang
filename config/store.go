package config

import (
	cmn "filestore-server/common"
)

const (
	// TempLocalRootDir : 本地临时存储地址的路径
	TempLocalRootDir = "./tmp/"
	// CurrentStoreType : 设置当前文件的存储类型
	CurrentStoreType = cmn.StoreS3
)