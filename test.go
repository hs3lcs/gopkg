package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/hs3lcs/gopkg/crypto"
	"github.com/hs3lcs/gopkg/dbms"
	"github.com/hs3lcs/gopkg/restapi"
	"github.com/joho/godotenv"
)

func init() {
	envFile := flag.String("e", ".env", "env file")
	flag.Parse()
	err := godotenv.Load(*envFile)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	dbmsTest()
	cryptoTest()
	// restapi
	fmt.Println("- restapi -")
	res, err := restapi.Call(restapi.ApiPack{
		Url:    "https://wkrh.info/api/v1",
		Method: "GET",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(res))
}

func dbmsTest() {
	db, err := dbms.Init(dbms.Config{
		RWDB: dbms.DBCON{
			HOST: os.Getenv("MYSQL_RW_HOST"),
			USER: os.Getenv("MYSQL_RW_USER"),
			PASS: os.Getenv("MYSQL_RW_PASS"),
		},
		RODB: dbms.DBCON{
			HOST: os.Getenv("MYSQL_RO_HOST"),
			USER: os.Getenv("MYSQL_RO_USER"),
			PASS: os.Getenv("MYSQL_RO_PASS"),
		},
		REDIS: dbms.RDBCON{
			DB:   1,
			HOST: os.Getenv("REDIS_HOST"),
			PASS: os.Getenv("REDIS_PASS"),
		},
	})
	if err != nil {
		fmt.Println("dbms init", err)
		return
	}
	// dbms mysql
	fmt.Println("- dbms mysql -")
	rows, err := db.RO.Query("SHOW DATABASES")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			fmt.Println(err)
		}
		fmt.Println(dbName)
	}
	// dbms cache
	fmt.Println("- dbms cache -")
	err = db.SetCache("test", "test cache", time.Minute*1)
	if err != nil {
		fmt.Println(err)
		return
	}
	str, err := db.GetCache("test")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("test:", str)
	all, err := db.GetAllCache()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("all:", all)
}

func cryptoTest() {
	jwtExp, _ := strconv.Atoi(os.Getenv("JWT_EXP"))
	crypto.Config = &crypto.Cfg{
		JWT_EXP: jwtExp,
		JWT_KEY: os.Getenv("JWT_KEY"),
	}
	fmt.Println("- crypto -")
	token, err := crypto.JwtEncode(crypto.JwtPack{})
	if err != nil {
		fmt.Println(err)
		return
	}
	jwt, err := crypto.JwtDecode(token)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", jwt)
	fmt.Printf("%+v\n", crypto.JwtParse(token))
	fmt.Println(crypto.HashString(32))
}
