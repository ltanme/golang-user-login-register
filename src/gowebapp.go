package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"fmt"

	"io"

	_ "github.com/go-sql-driver/mysql"
)

//全局变量声明 mysql查询
var db = &sql.DB{}

//全局变量声明 模版位置
var template_base_path = "/Users/winfan/Documents/workspace/golang-user-login-register/src/template/"

//静态目录设置
var static_access_path = "/Users/winfan/Documents/workspace/golang-user-login-register/src/static"

func init() {
	var err error
	//sql.open不是真真意义上打开mysql连接，当只有执行exec查询。。。操作时才去连接mysql
	db, err = sql.Open("mysql", "root:123456@/test?charset=utf8")
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	db.Ping()
	checkErr(err)
}

var tdata = map[string]interface{}{
	"title": `>title`,
	"body":  `/admin/body.html`,
	"js":    `/admin/js.tmpl`,
	"href":  ">>>",
	"name":  "admin",
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

/*
设置模版
*/
func setTemplate(sname string, w http.ResponseWriter, data map[string]interface{}) {
	t, _ := template.ParseFiles(template_base_path + sname)
	t.Execute(w, data)
}

/*
登录页面
*/
func loginPage(w http.ResponseWriter, r *http.Request) {
	log.Println("open file error555 !")
	setTemplate("login.html", w, tdata)
}

/*
登录处理
*/
func loginHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fmt.Println(r.Form["email"][0])
	fmt.Println(r.Form["password"][0])
	rows, err := db.Query("SELECT username,email FROM userinfo WHERE email=? and password=?", r.Form["email"][0], r.Form["password"][0])
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var isUser = false

	for rows.Next() {
		var username string
		var email string
		if err := rows.Scan(&username, &email); err != nil {
			log.Fatal(err)
		}
		fmt.Println(username, email)
		isUser = true
	}
	if isUser {
		http.Redirect(w, r, "./admin", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "./login", http.StatusSeeOther)
	}

}

/*
注册页面
*/
func registerPage(w http.ResponseWriter, r *http.Request) {
	setTemplate("register.html", w, tdata)
}

/*
注册处理
*/
func registerHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tx, _ := db.Begin()
	tx.Exec("INSERT INTO userinfo(username,email,password) values(?,?,?)",
		r.Form["username"][0],
		r.Form["email"][0],
		r.Form["password"][0])
	tx.Commit()
	//log.Println(err)
	http.Redirect(w, r, "./login", http.StatusSeeOther)
}

/*
后台首页页面
*/
func adminPage(w http.ResponseWriter, r *http.Request) {
	setTemplate("admin.html", w, tdata)
}
func api(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "{asfdsaf}")
}
func main() {

	fmt.Println("开始启动！！")
	h := http.FileServer(http.Dir(static_access_path)) //设置webapp的静态目录
	http.Handle("/static/", http.StripPrefix("/static/", h))

	//定义路由
	http.HandleFunc("/admin", adminPage)
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/login-handler", loginHandle)
	http.HandleFunc("/register", registerPage)
	http.HandleFunc("/register-handler", registerHandle)
	http.HandleFunc("/api", api)
	//绑定端口并监听
	serr := http.ListenAndServe(":9090", nil)
	if serr != nil {
		log.Fatal("监听端口", serr)
	}

}
