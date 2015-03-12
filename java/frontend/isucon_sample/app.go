package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/draftcode/isucon_misc/grizzly"
	"github.com/draftcode/isucon_practice/isucon3_qual/sessions"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	maxConnectionCount = 256
	memosPerPage       = 100
	listenAddr         = ":80"
	sessionName        = "isucon_session"
	tmpDir             = "/tmp/"
	markdownCommand    = "../bin/markdown"
	dbConnPoolSize     = 10
	memcachedServer    = "localhost:11212"
	sessionSecret      = "kH<{11qpic*gf0e21YK7YtwyUvE9l<1r>yX8R-Op"
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

type User struct {
	Id         int
	Username   string
	Password   string
	Salt       string
	LastAccess string
}

type Memo struct {
	Id        int
	User      int
	Title     string
	Content   string
	IsPrivate int
	CreatedAt string
	UpdatedAt string
	Username  string
}

type Memos []*Memo

type View struct {
	User      *User
	Memo      *Memo
	Memos     *Memos
	Page      int
	PageStart int
	PageEnd   int
	Total     int
	Older     *Memo
	Newer     *Memo
	Token     string
}

var (
	dbConnPool chan *sql.DB
	baseUrl    *url.URL
	fmap       = template.FuncMap{
		"gen_markdown": func(s string) template.HTML {
			return template.HTML(blackfriday.MarkdownBasic([]byte(s)))
		},
	}
	tmpl      = template.Must(template.New("tmpl").Funcs(fmap).ParseGlob("templates/*.html"))
	redisPool *redis.Pool
)

func getToken(session *sessions.Session) string {
	if session == nil || session.Values["token"] == nil {
		return ""
	} else {
		return session.Values["token"].(string)
	}
}

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

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	env := os.Getenv("ISUCON_ENV")
	if env == "" {
		env = "local"
	}
	// config := loadConfig("/home/isucon/webapp/config/" + env + ".json")
	db := config.Database
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8",
		db.Username, db.Password, db.Host, db.Port, db.Dbname,
	)
	log.Printf("db: %s", connectionString)

	dbConnPool = make(chan *sql.DB, dbConnPoolSize)
	for i := 0; i < dbConnPoolSize; i++ {
		conn, err := sql.Open("mysql", connectionString)
		if err != nil {
			log.Panicf("Error opening database: %v", err)
		}
		dbConnPool <- conn
		defer conn.Close()
	}

	redisPool = newPool(":6379")

	r := mux.NewRouter()
	r.HandleFunc("/", grizzly.WrapHandleFunc("Top", topHandler))
	r.HandleFunc("/signin", grizzly.WrapHandleFunc("SignIn", signinHandler)).Methods("GET", "HEAD")
	r.HandleFunc("/signin", grizzly.WrapHandleFunc("SignInPost", signinPostHandler)).Methods("POST")
	r.HandleFunc("/signout", grizzly.WrapHandleFunc("SignOut", signoutHandler))
	r.HandleFunc("/mypage", grizzly.WrapHandleFunc("MyPage", mypageHandler))
	r.HandleFunc("/memo/{memo_id}", grizzly.WrapHandleFunc("Memo", memoHandler)).Methods("GET", "HEAD")
	r.HandleFunc("/memo", grizzly.WrapHandleFunc("MemoPost", memoPostHandler)).Methods("POST")
	r.HandleFunc("/recent/{page:[0-9]+}", grizzly.WrapHandleFunc("Recent", recentHandler))
	r.PathPrefix("/").Handler(grizzly.WrapHandler("Static", http.FileServer(http.Dir("./public/"))))
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

//func loadConfig(filename string) *Config {
//	log.Printf("loading config file: %s", filename)
//	f, err := ioutil.ReadFile(filename)
//	if err != nil {
//		log.Fatal(err)
//		os.Exit(1)
//	}
//	var config Config
//	err = json.Unmarshal(f, &config)
//	if err != nil {
//		log.Fatal(err)
//		os.Exit(1)
//	}
//	return &config
//}

func prepareHandler(w http.ResponseWriter, r *http.Request) {
	if h := r.Header.Get("X-Forwarded-Host"); h != "" {
		baseUrl, _ = url.Parse("http://" + h)
	} else {
		baseUrl, _ = url.Parse("http://" + r.Host)
	}
}

func loadSession(w http.ResponseWriter, r *http.Request) (session *sessions.Session, err error) {
	store := sessions.NewMemcacheStore(memcachedServer, []byte(sessionSecret))
	return store.Get(r, sessionName)
}

func getUser(w http.ResponseWriter, r *http.Request, dbConn *sql.DB, session *sessions.Session) *User {
	userId := session.Values["user_id"]
	if userId == nil {
		return nil
	}
	user := &User{}
	rows, err := dbConn.Query("SELECT * FROM users WHERE id=?", userId)
	if err != nil {
		serverError(w, err)
		return nil
	}
	if rows.Next() {
		rows.Scan(&user.Id, &user.Username, &user.Password, &user.Salt, &user.LastAccess)
		rows.Close()
	}
	if user != nil {
		w.Header().Add("Cache-Control", "private")
	}
	return user
}

func antiCSRF(w http.ResponseWriter, r *http.Request, session *sessions.Session) bool {
	if r.FormValue("sid") != session.Values["token"] {
		code := http.StatusBadRequest
		http.Error(w, http.StatusText(code), code)
		return true
	}
	return false
}

func serverError(w http.ResponseWriter, err error) {
	log.Printf("error: %s", err)
	code := http.StatusInternalServerError
	http.Error(w, http.StatusText(code), code)
}

func notFound(w http.ResponseWriter) {
	code := http.StatusNotFound
	http.Error(w, http.StatusText(code), code)
}

func topHandler(w http.ResponseWriter, r *http.Request) {
	session, err := loadSession(w, r)
	if err != nil {
		serverError(w, err)
		return
	}
	prepareHandler(w, r)
	dbConn := <-dbConnPool
	defer func() {
		dbConnPool <- dbConn
	}()
	user := getUser(w, r, dbConn, session)

	redisConn := redisPool.Get()
	defer redisConn.Close()

	var totalCount uint64
	totalCount, err = redis.Uint64(redisConn.Do("zcard", "public_recent"))
	if err != nil {
		serverError(w, err)
		return
	}

	memos := make(Memos, 0)
	recentIds, err := redis.Strings(redisConn.Do("zrevrange", "public_recent", 0, memosPerPage-1))
	if err != nil {
		serverError(w, err)
		return
	}
	if len(recentIds) > 0 {
		for i, recentId := range recentIds {
			recentIds[i] = "memos.id = " + recentId
		}
		rows, err := dbConn.Query(`
			SELECT memos.id, memos.title, memos.created_at, users.username
			FROM memos
			INNER JOIN users ON memos.user = users.id
			WHERE ` + strings.Join(recentIds, " OR ") + `
			ORDER BY memos.created_at DESC, memos.id DESC`,
		)
		if err != nil {
			serverError(w, err)
			return
		}
		for rows.Next() {
			memo := Memo{}
			rows.Scan(&memo.Id, &memo.Title, &memo.CreatedAt, &memo.Username)
			memos = append(memos, &memo)
		}
		rows.Close()
	}

	v := &View{
		Total:     int(totalCount),
		Page:      0,
		PageStart: 1,
		PageEnd:   memosPerPage,
		Memos:     &memos,
		User:      user,
		Token:     getToken(session),
	}
	if err = tmpl.ExecuteTemplate(w, "index", v); err != nil {
		serverError(w, err)
	}
}

func recentHandler(w http.ResponseWriter, r *http.Request) {
	session, err := loadSession(w, r)
	if err != nil {
		serverError(w, err)
		return
	}
	prepareHandler(w, r)
	dbConn := <-dbConnPool
	defer func() {
		dbConnPool <- dbConn
	}()
	user := getUser(w, r, dbConn, session)
	vars := mux.Vars(r)
	page, _ := strconv.Atoi(vars["page"])

	redisConn := redisPool.Get()
	defer redisConn.Close()

	b := grizzly.StartCodeBlock("recent.total")
	var totalCount uint64
	totalCount, err = redis.Uint64(redisConn.Do("zcard", "public_recent"))
	if err != nil {
		serverError(w, err)
		return
	}
	b.Close()

	b = grizzly.StartCodeBlock("recent.rows")
	memos := make(Memos, 0)
	recentIds, err := redis.Strings(redisConn.Do("zrevrange", "public_recent", memosPerPage*page, memosPerPage*(page+1)-1))
	if err != nil {
		serverError(w, err)
		return
	}
	if len(recentIds) > 0 {
		for i, recentId := range recentIds {
			recentIds[i] = "memos.id = " + recentId
		}
		rows, err := dbConn.Query(`
			SELECT memos.id, memos.title, memos.created_at, users.username
			FROM memos
			INNER JOIN users ON memos.user = users.id
			WHERE ` + strings.Join(recentIds, " OR ") + `
			ORDER BY memos.created_at DESC, memos.id DESC`,
		)
		if err != nil {
			serverError(w, err)
			return
		}
		for rows.Next() {
			memo := Memo{}
			rows.Scan(&memo.Id, &memo.Title, &memo.CreatedAt, &memo.Username)
			memos = append(memos, &memo)
		}
		rows.Close()
	} else {
		notFound(w)
		return
	}
	b.Close()

	v := &View{
		Total:     int(totalCount),
		Page:      page,
		PageStart: memosPerPage*page + 1,
		PageEnd:   memosPerPage * (page + 1),
		Memos:     &memos,
		User:      user,
		Token:     getToken(session),
	}

	b = grizzly.StartCodeBlock("recent.tmpl")
	if err = tmpl.ExecuteTemplate(w, "index", v); err != nil {
		serverError(w, err)
	}
	b.Close()
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	session, err := loadSession(w, r)
	if err != nil {
		serverError(w, err)
		return
	}
	prepareHandler(w, r)
	dbConn := <-dbConnPool
	defer func() {
		dbConnPool <- dbConn
	}()
	user := getUser(w, r, dbConn, session)

	v := &View{
		User:  user,
		Token: getToken(session),
	}
	if err := tmpl.ExecuteTemplate(w, "signin", v); err != nil {
		serverError(w, err)
		return
	}
}

func signinPostHandler(w http.ResponseWriter, r *http.Request) {
	session, err := loadSession(w, r)
	if err != nil {
		serverError(w, err)
		return
	}
	prepareHandler(w, r)
	dbConn := <-dbConnPool
	defer func() {
		dbConnPool <- dbConn
	}()

	username := r.FormValue("username")
	password := r.FormValue("password")
	user := &User{}
	rows, err := dbConn.Query("SELECT id, username, password, salt FROM users WHERE username=?", username)
	if err != nil {
		serverError(w, err)
		return
	}
	if rows.Next() {
		rows.Scan(&user.Id, &user.Username, &user.Password, &user.Salt)
	}
	rows.Close()
	if user.Id > 0 {
		h := sha256.New()
		h.Write([]byte(user.Salt + password))
		if user.Password == fmt.Sprintf("%x", h.Sum(nil)) {
			session.Values["user_id"] = user.Id
			session.Values["token"] = fmt.Sprintf("%x", securecookie.GenerateRandomKey(32))
			if err := session.Save(r, w); err != nil {
				serverError(w, err)
				return
			}
			if _, err := dbConn.Exec("UPDATE users SET last_access=now() WHERE id=?", user.Id); err != nil {
				serverError(w, err)
				return
			} else {
				http.Redirect(w, r, "/mypage", http.StatusFound)
			}
			return
		}
	}
	v := &View{
		Token: getToken(session),
	}
	if err := tmpl.ExecuteTemplate(w, "signin", v); err != nil {
		serverError(w, err)
		return
	}
}

func signoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := loadSession(w, r)
	if err != nil {
		serverError(w, err)
		return
	}
	prepareHandler(w, r)
	if antiCSRF(w, r, session) {
		return
	}

	http.SetCookie(w, sessions.NewCookie(sessionName, "", &sessions.Options{MaxAge: -1}))
	http.Redirect(w, r, "/", http.StatusFound)
}

func mypageHandler(w http.ResponseWriter, r *http.Request) {
	session, err := loadSession(w, r)
	if err != nil {
		serverError(w, err)
		return
	}
	prepareHandler(w, r)
	dbConn := <-dbConnPool
	defer func() {
		dbConnPool <- dbConn
	}()

	user := getUser(w, r, dbConn, session)
	if user == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	rows, err := dbConn.Query("SELECT id, title, is_private, created_at, updated_at FROM memos WHERE user=? ORDER BY created_at DESC", user.Id)
	if err != nil {
		serverError(w, err)
		return
	}
	memos := make(Memos, 0)
	for rows.Next() {
		memo := Memo{}
		rows.Scan(&memo.Id, &memo.Title, &memo.IsPrivate, &memo.CreatedAt, &memo.UpdatedAt)
		memos = append(memos, &memo)
	}
	v := &View{
		Memos: &memos,
		User:  user,
		Token: getToken(session),
	}
	if err = tmpl.ExecuteTemplate(w, "mypage", v); err != nil {
		serverError(w, err)
	}
}

func memoHandler(w http.ResponseWriter, r *http.Request) {
	session, err := loadSession(w, r)
	if err != nil {
		serverError(w, err)
		return
	}
	prepareHandler(w, r)
	vars := mux.Vars(r)
	memoId := vars["memo_id"]
	dbConn := <-dbConnPool
	defer func() {
		dbConnPool <- dbConn
	}()
	user := getUser(w, r, dbConn, session)

	rows, err := dbConn.Query("SELECT id, user, content, is_private, created_at, updated_at FROM memos WHERE id=?", memoId)
	if err != nil {
		serverError(w, err)
		return
	}
	memo := &Memo{}
	if rows.Next() {
		rows.Scan(&memo.Id, &memo.User, &memo.Content, &memo.IsPrivate, &memo.CreatedAt, &memo.UpdatedAt)
		rows.Close()
	} else {
		notFound(w)
		return
	}
	if memo.IsPrivate == 1 {
		if user == nil || user.Id != memo.User {
			notFound(w)
			return
		}
	}
	rows, err = dbConn.Query("SELECT username FROM users WHERE id=?", memo.User)
	if err != nil {
		serverError(w, err)
		return
	}
	if rows.Next() {
		rows.Scan(&memo.Username)
		rows.Close()
	}

	var cond string
	if user != nil && user.Id == memo.User {
		cond = ""
	} else {
		cond = "AND is_private=0"
	}
	rows, err = dbConn.Query("SELECT id FROM memos WHERE user=? "+cond+" ORDER BY created_at", memo.User)
	if err != nil {
		serverError(w, err)
		return
	}
	memos := make(Memos, 0)
	for rows.Next() {
		m := Memo{}
		rows.Scan(&m.Id)
		memos = append(memos, &m)
	}
	rows.Close()
	var older *Memo
	var newer *Memo
	for i, m := range memos {
		if m.Id == memo.Id {
			if i > 0 {
				older = memos[i-1]
			}
			if i < len(memos)-1 {
				newer = memos[i+1]
			}
		}
	}

	v := &View{
		User:  user,
		Memo:  memo,
		Older: older,
		Newer: newer,
		Token: getToken(session),
	}
	if err = tmpl.ExecuteTemplate(w, "memo", v); err != nil {
		serverError(w, err)
	}
}

func memoPostHandler(w http.ResponseWriter, r *http.Request) {
	session, err := loadSession(w, r)
	if err != nil {
		serverError(w, err)
		return
	}
	prepareHandler(w, r)
	if antiCSRF(w, r, session) {
		return
	}
	dbConn := <-dbConnPool
	defer func() {
		dbConnPool <- dbConn
	}()

	user := getUser(w, r, dbConn, session)
	if user == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	var isPrivate int
	if r.FormValue("is_private") == "1" {
		isPrivate = 1
	} else {
		isPrivate = 0
	}
	content := r.FormValue("content")
	title := strings.Split(content, "\n")[0]
	result, err := dbConn.Exec(
		"INSERT INTO memos (user, title, content, is_private, created_at) VALUES (?, ?, ?, ?, now())",
		user.Id, title, content, isPrivate,
	)
	if err != nil {
		serverError(w, err)
		return
	}

	newId, _ := result.LastInsertId()
	if isPrivate == 0 {
		redisConn := redisPool.Get()
		if _, err := redisConn.Do("zadd", "public_recent", newId, newId); err != nil {
			serverError(w, err)
			return
		}
		redisConn.Close()
	}

	http.Redirect(w, r, fmt.Sprintf("/memo/%d", newId), http.StatusFound)
}
