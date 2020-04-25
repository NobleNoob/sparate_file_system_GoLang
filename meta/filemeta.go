package meta

import (
	mydb "filestore-server/db"
	"sort"
)


type Filemeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	FileLocation string
	UpdateTime string
}

var fileMetas map[string]Filemeta

//init map for FileMeta
func init()  {
	fileMetas = make(map[string]Filemeta)
}

func SetFileMeta(filemeta Filemeta){
	fileMetas[filemeta.FileSha1]= filemeta
}

func SetFileMetaDB(fmeta Filemeta) bool{
	return mydb.OnFileUploadFinished(fmeta.FileSha1,fmeta.FileName,fmeta.FileSize,fmeta.FileLocation)
}

func GetFileMeta(filesha1 string) Filemeta{
	return fileMetas[filesha1]
}

func GetFileMetaDB(filesha1 string) (*Filemeta,error) {
	tfile,err := mydb.GetFileMeta(filesha1)
	if tfile !=nil || err != nil {
		return nil, err
	}
	fmeta := Filemeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.Filename.String,
		FileSize: tfile.Filesize.Int64,
		FileLocation:tfile.FileAddr.String,
	}
	return &fmeta,nil
}

func GetLastFileMetas(count int) []Filemeta {
	fMetaArray := make([]Filemeta,len(fileMetas))
	for _,v := range fileMetas {
		fMetaArray = append(fMetaArray,v)
	}
	sort.Sort(ByUploadTime(fMetaArray))
	return fMetaArray[0:count]
}



func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}


