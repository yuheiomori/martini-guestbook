package main

import (
	"html/template"
	"net/http"
	"time"

	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/martini-contrib/render"
)

type Greeting struct {
	Id       int64
	Name     string `sql:"size:100"`
	Comment  string `sql:"size:200"`
	CreateAt int64
}

func main() {
	m := martini.Classic()
	// 静的ファイルの利用
	m.Use(martini.Static("static"))

	// テンプレート関数の登録
	m.Use(render.Renderer(render.Options{
		Directory:  "templates",                // Specify what path to load the templates from.
		Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
		Funcs: []template.FuncMap{{
			"nl2br":      nl2br,
			"htmlquote":  htmlquote,
			"str2html":   str2html,
			"dateformat": dateformat,
		}},
	}))

	// データベース初期化
	dbmap := initDb()
	defer dbmap.Db.Close()

	// トップページ
	m.Get("/", func(render render.Render) {
		var greetings []Greeting
		_, err := dbmap.Select(&greetings, "select * from greetings order by CreateAt desc")
		if err != nil {
			panic(err)
		}
		render.HTML(200, "index", greetings)
	})

	// 投稿
	m.Post("/post", func(w http.ResponseWriter, r *http.Request, render render.Render) {
		greeting := Greeting{
			Name:     r.FormValue("name"),
			Comment:  r.FormValue("comment"),
			CreateAt: time.Now().UnixNano(),
		}
		dbmap.Insert(&greeting)
		render.Redirect("/", 302)
	})
	m.Run()
}
