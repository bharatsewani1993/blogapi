package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

// Article stores all article data
type Article struct {
	ID           int    `json:"id"`
	Name         string `json:"nickname"`
	Title        string `json:"title"`
	Creationdate string `json:"creation_date"`
}

//Sarticle stores a single article
type Sarticle struct {
	ID           int    `json:"id"`
	Name         string `json:"nickname"`
	Title        string `json:"title"`
	Content      string `json:"Content"`
	Creationdate string `json:"Creationdate"`
}

// Comment stores comments
type Comment struct {
	ID            int    `json:"id"`
	Ccontent      string `json:"c_content"`
	CNickname     string `json:"c_nickname"`
	CCreationdate string `json:"c_creation_date"`
	Scontent      string `json:"s_content"`
	SNickName     string `json:"s_nickname"`
	SCreationdate string `json:"s_creation_date"`
}

// Subcomment stores comment reply
type Subcomment struct {
	ID           int    `json:"id"`
	CommentID    int    `json:"comment_id"`
	Nickname     string `json:"nickname"`
	Content      string `json:"content"`
	Creationdate string `json:"creation_date"`
	Active       string `json:"active"`
}

//initialize database connection
func dbconnection() *sql.DB {
	db, err := sql.Open("mysql", "root:P3NT3ST3R@/golangtask")
	if err != nil {
		fmt.Printf("Error %v\n", err)
	}
	return db
}

//load the index page for documentation
func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome to blog api page!\n")
	fmt.Fprintf(w, "\n load list of articles page wise:  /page/id\n/page/1 \n")
	fmt.Fprintf(w, "\n load article and its content: /article/articleid\n/article/1 \n")
	fmt.Fprintf(w, "\n load comment and its sub comments: /comment/commentid\n/comment/1 \n")
	fmt.Fprintf(w, "\n\n\n Post links are mentioned in image attachment \n")
}

//list out all the articles
func loadpage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db := dbconnection()

	pageid := ps.ByName("id")
	s, err := strconv.Atoi(pageid)
	//to load 20 article per page
	var Ipageid = 0
	if s > 1 {
		Ipageid = s * 10
	}
	fmt.Println(Ipageid)
	stmt, err := db.Prepare("SELECT id,nickname,title,creation_date from articles where active=1 limit ?,20")
	if err != nil {
		fmt.Printf("Error %v\n", err)
	}

	results, err := stmt.Query(Ipageid)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	for results.Next() {
		var article Article
		err = results.Scan(&article.ID, &article.Name, &article.Title, &article.Creationdate)
		if err != nil {
			panic(err.Error())
		}
		ar, err := json.Marshal(article)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(ar))
		fmt.Fprintf(w, string(ar))
	}
}

//load specific article
func loadarticle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db := dbconnection()
	articleid := ps.ByName("id")

	stmt, err := db.Prepare("SELECT id,nickname,title,content,creation_date from articles where id=? and active=1")
	if err != nil {
		fmt.Printf("Error %v\n", err)
	}

	results, err := stmt.Query(articleid)
	if err != nil {
		fmt.Printf("Error %v\n", err)
	}

	for results.Next() {
		var sarticle Sarticle
		err = results.Scan(&sarticle.ID, &sarticle.Name, &sarticle.Title, &sarticle.Content, &sarticle.Creationdate)
		if err != nil {
			panic(err.Error())
		}
		sar, err := json.Marshal(sarticle)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(sar))
		fmt.Fprintf(w, string(sar))
	}

}

//load comments for an article
func loadcomment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db := dbconnection()
	commentid := ps.ByName("id")

	stmt, err := db.Prepare("SELECT C.id,C.content as c_content,C.nickname as c_nickname,C.creation_date as c_creation_date, S.content as s_content,S.nickname as s_nickname,S.creation_date as s_creation_date from comments C JOIN sub_comments S ON C.id=S.comment_id where C.article_id=? and C.active =1")
	if err != nil {
		fmt.Printf("Error %v\n", err)
	}

	results, err := stmt.Query(commentid)
	if err != nil {
		fmt.Printf("Error %v\n", err)
	}

	for results.Next() {
		var comment Comment
		err = results.Scan(&comment.ID, &comment.Ccontent, &comment.CNickname, &comment.CCreationdate, &comment.Scontent, &comment.SNickName, &comment.SCreationdate)
		if err != nil {
			panic(err.Error())
		}
		cm, err := json.Marshal(comment)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(cm))
		fmt.Fprintf(w, string(cm))
	}
}

//post an article
func postarticle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()

	//required field validations
	if len(r.Form.Get("nickname")) != 0 && len(r.Form.Get("title")) != 0 && len(r.Form.Get("content")) != 0 {

		db := dbconnection()
		currenttime := time.Now().Local()
		ct := currenttime.Format("2006-01-02")
		stmt, err := db.Prepare("INSERT articles SET nickname=?,title=?,content=?,creation_date=?")
		checkErr(err)
		_, err = stmt.Exec(r.Form["nickname"][0], r.Form["title"][0], r.Form["content"][0], ct)
		checkErr(err)
		fmt.Fprintf(w, "Article Successfuly Created")
	} else {
		fmt.Printf("Please send all required fields")
		fmt.Fprintf(w, "Please send all required fields")
	}
}

//post comment for article
func postcomment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	//required fields validation
	if len(r.Form.Get("nickname")) != 0 && len(r.Form.Get("article_id")) != 0 && len(r.Form.Get("content")) != 0 {
		db := dbconnection()
		currenttime := time.Now().Local()
		ct := currenttime.Format("2006-01-02")
		stmt, err := db.Prepare("INSERT comments SET article_id=?,nickname=?,content=?,creation_date=?")
		checkErr(err)
		_, err = stmt.Exec(r.Form["article_id"][0], r.Form["nickname"][0], r.Form["content"][0], ct)
		checkErr(err)
		fmt.Fprintf(w, "Comment Successfuly Created")
	} else {
		fmt.Printf("Please send all required fields")
		fmt.Fprintf(w, "Please send all required fields")
	}
}

//post sub comments for an article
func postsubcomment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	//required field validation
	if len(r.Form.Get("nickname")) != 0 && len(r.Form.Get("comment_id")) != 0 && len(r.Form.Get("content")) != 0 {
		db := dbconnection()
		currenttime := time.Now().Local()
		ct := currenttime.Format("2006-01-02")
		stmt, err := db.Prepare("INSERT sub_comments SET comment_id=?, nickname=?, content=?, creation_date=?")
		checkErr(err)
		_, err = stmt.Exec(r.Form["comment_id"][0], r.Form["nickname"][0], r.Form["content"][0], ct)
		checkErr(err)
		fmt.Fprintf(w, "Subcomment Successfuly Created")
	} else {
		fmt.Printf("Please send all required fields")
		fmt.Fprintf(w, "Please send all required fields")
	}
}

//check error function
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	router := httprouter.New()
	router.GET("/", index)
	router.GET("/page/:id", loadpage)
	router.GET("/article/:id", loadarticle)
	router.GET("/comment/:id", loadcomment)
	router.POST("/postarticle/", postarticle)
	router.POST("/postcomment/", postcomment)
	router.POST("/postsubcomment", postsubcomment)
	log.Fatal(http.ListenAndServe(":8080", router))
}
