package mq

import (
	cmn "filestore-server/common"
)

type TransferData struct {
	FileHash      string
	// CurLocation 本地临时存储路径
	CurLocation   string
	// DestLocation 目的地存储路径
	DestLocation  string
	DestStoreType cmn.StoreType
}
