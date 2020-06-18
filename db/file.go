package db
import (
	"database/sql"
	mydb "filestore-server/db/mysql"
	"fmt"
)

func OnFileUploadFinished(filehash string,filename string,filesize int64, fileaddr string) bool {
	stmt,err := mydb.DBconn().Prepare(
		"insert ignore into tbl_file(`file_sha1`,`file_name`,`file_size`,`file_addr`,`status`) values (?,?,?,?,1)")
	if err !=  nil {
		fmt.Printf("failed on Prepare step" +err.Error())
	}

	defer stmt.Close()
	result,err := stmt.Exec(filehash,filename,filesize,fileaddr)
	if err != nil {
		fmt.Printf(err.Error())
		return false
	}

	if rf,err := result.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("File with hash: %s has been uploaded before", filehash)
		}
		return true
	}
	return false
}


type TableFile struct {
	FileHash string
	Filename sql.NullString
	Filesize sql.NullInt64
	FileAddr sql.NullString
}

func GetFileMeta(filehash string) (*TableFile ,error){
	rel,err := mydb.DBconn().Prepare("select file_sha1,file_addr,file_name,file_size from tbl_file where file_sha1=? and status=1 limit 1")

	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	defer rel.Close()

	tfile := TableFile{}
	err = rel.QueryRow(filehash).Scan(&tfile.FileHash,&tfile.FileAddr,&tfile.Filename,&tfile.Filesize)
	if err != nil {
		fmt.Printf(err.Error())
		return nil,err
	}
	return &tfile,nil
}