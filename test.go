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
	// config
	jwtExp, _ := strconv.Atoi(os.Getenv("JWT_EXP"))
	crypto.Config = &crypto.Cfg{
		JWT_EXP: jwtExp,
		JWT_KEY: os.Getenv("JWT_KEY"),
	}
	dbms.Config = &dbms.Cfg{
		MYSQL_RW_HOST: os.Getenv("MYSQL_RW_HOST"),
		MYSQL_RW_USER: os.Getenv("MYSQL_RW_USER"),
		MYSQL_RW_PASS: os.Getenv("MYSQL_RW_PASS"),
		MYSQL_RO_HOST: os.Getenv("MYSQL_RO_HOST"),
		MYSQL_RO_USER: os.Getenv("MYSQL_RO_USER"),
		MYSQL_RO_PASS: os.Getenv("MYSQL_RO_PASS"),
		REDIS_HOST:    os.Getenv("REDIS_HOST"),
		REDIS_PASS:    os.Getenv("REDIS_PASS"),
	}
	// dbms mysql
	fmt.Println("- dbms mysql -")
	db, err := dbms.Connect("RO")
	if err != nil {
		fmt.Println(err)
		return
	}
	rows, err := db.Query("SHOW DATABASES")
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
	// dbms redis
	fmt.Println("- dbms redis -")
	err = dbms.CacheSet("cache", "test cache", time.Minute*1)
	if err != nil {
		fmt.Println(err)
		return
	}
	str, err := dbms.CacheGet("cache")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("cache:", str)
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
	// crypto
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
