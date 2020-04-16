package db

import (
	mydb "filestore-server/db/mysql"
	"fmt"
)
//User model
type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

type User_token struct {
	Username	string
	Token 		string
}

func UserSignUp(username string,password string) bool{
	stmt,err := mydb.DBconn().Prepare("insert ignore into tbl_user (`user_name`,`user_pwd`) values (?,?)")
	if err != nil {
		fmt.Printf("Failed to insert,error:" + err.Error())
		return false
	}
	defer stmt.Close()
	ret,err:=stmt.Exec(username,password)
	if err != nil {
		fmt.Printf("Failed to insert,error:" + err.Error())
		return false
	}
	if rowsAffected,err := ret.RowsAffected(); nil == err && rowsAffected >0 {
		return true
	}
	return false
}

func UserSignin(username string,encpwd string) bool {

	stmt,err:=mydb.DBconn().Prepare("select * from tbl_user where user_name =? limit 1")
	if err != nil {
		fmt.Printf(err.Error())
		return false
	}
	rows,err := stmt.Query(username)
	if err != nil {
		fmt.Printf(err.Error())
		return false
	} else if rows== nil {
		fmt.Printf("username not found:" + username)
		return false
	}

	pRows := mydb.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}
	return false
}

func UpdateToken(username string, token string) bool {
	stmt,err :=mydb.DBconn().Prepare(
		"replace into tbl_user_token(`user_name`,`user_token`)  values (?,?)")
	if err != nil {
		fmt.Printf(err.Error())
		return false
	}
	defer stmt.Close()
	_,err = stmt.Exec(username ,token)
	if err != nil {
		fmt.Printf(err.Error())
		return false
	}
	return true
}

func GetUserInfo(username string) (User, error) {
	user := User{}

	stmt, err := mydb.DBconn().Prepare(
		"select user_name,signup_at from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	defer stmt.Close()

	// 执行查询的操作
	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		return user, err
	}
	return user, nil
}

func GetUserToken(username string) (User_token,error) {
	user_token:=User_token{}
	stmt, err := mydb.DBconn().Prepare(
		"select user_name,user_token from tbl_user_token where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return user_token,err
	}
	defer stmt.Close()
	err = stmt.QueryRow(username).Scan(&user_token.Token)
	if err != nil {
		return user_token, err
	}
	return user_token, nil
}