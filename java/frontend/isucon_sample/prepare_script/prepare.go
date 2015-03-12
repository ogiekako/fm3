package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Config struct {
	Database struct {
		Dbname   string `json:"dbname"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"database"`
}

var (
	dbConn    *sql.DB
	redisPool *redis.Pool
)

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     8,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}


func loadConfig(filename string) *Config {
	log.Printf("loading config file: %s", filename)
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	var config Config
	err = json.Unmarshal(f, &config)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return &config
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	env := os.Getenv("ISUCON_ENV")
	if env == "" {
		env = "local"
	}
	config := loadConfig("/home/isucon/webapp/config/" + env + ".json")
	db := config.Database
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8",
		db.Username, db.Password, db.Host, db.Port, db.Dbname,
	)
	log.Printf("db: %s", connectionString)

	var err error

	dbConn, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	redisPool = newPool(":6379")
	redisConn := redisPool.Get()
	redisConn.Do("flushdb")
	redisConn.Close()

	updateMemos()
}

type Memo struct {
	Id        int
	Content   string
	IsPrivate int
}

func updateMemos() {
	var wg sync.WaitGroup
	c := make(chan *Memo)
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			redisConn := redisPool.Get()
			for memo := range c {
				title := strings.Split(memo.Content, "\n")[0]
				dbConn.Exec("UPDATE memos SET title=? WHERE id=?", title, memo.Id)
				if memo.IsPrivate == 0 {
					redisConn.Do("zadd", "public_recent", memo.Id, memo.Id)
				}
			}
			redisConn.Close()
			wg.Done()
		}()
	}
	rows, err := dbConn.Query("SELECT memos.id, memos.content, is_private FROM memos")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var id int
		var content string
		var isPrivate int
		rows.Scan(&id, &content, &isPrivate)
		c <- &Memo{id, content, isPrivate}
	}
	rows.Close()
	close(c)
	wg.Wait()
}
