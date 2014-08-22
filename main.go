package main

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
	"strings"
)

type DbContext struct {
	Dbmap *gorp.DbMap
}

type Member struct {
	Id         int
	First_Name string
	Last_Name  string
}

func NewDbContext() DbContext {
	db, err := sql.Open("postgres", "user=awhite password=password dbname=awhite sslmode=disable")
	if err != nil {

	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	return DbContext{Dbmap: dbmap}
}

func (self DbContext) CreateMember(first_name string, last_name string) {
	_, err := self.Dbmap.Exec("INSERT INTO Member (First_Name, Last_Name) VALUES ($1,$2)", first_name, last_name)
	if err != nil {
		panic(err)
	}
}

func (self DbContext) ReadMembers() []Member {
	var members []Member
	_, err := self.Dbmap.Select(&members, "SELECT Id, First_Name, Last_Name FROM member")
	if err != nil {
		panic(err)
	}
	return members
}

func (self DbContext) GetMember(id int) Member {
	var members []Member
	_, err := self.Dbmap.Select(&members, "SELECT Id, First_Name, Last_Name FROM member WHERE Id = $1 LIMIT 1", id)
	if err != nil {
		panic(err)
	}
	return members[0]
}

func (self DbContext) UpdateMember(id int, first_name string, last_name string) {
	_, err := self.Dbmap.Exec("UPDATE Member SET First_Name = $1, Last_Name = $2 WHERE Id = $3", first_name, last_name, id)
	if err != nil {
		panic(err)
	}
}

func (self DbContext) DeleteMember(id int) {
	_, err := self.Dbmap.Exec("DELETE FROM Member WHERE Id = $1", id)
	if err != nil {
		panic(err)
	}
}

func CloseDatabase(DbContext *DbContext) {
	DbContext.Dbmap.Db.Close()
}

func main() {
	r := gin.Default()
	r.LoadHTMLTemplates("templates/*")
	dbContext := NewDbContext()

	r.GET("/", func(c *gin.Context) {
		members := dbContext.ReadMembers()
		obj := gin.H{"members": members}
		c.HTML(200, "index.tmpl", obj)
	})

	r.POST("/create", func(c *gin.Context) {
		firstName := strings.TrimSpace(c.Request.FormValue("First_Name"))
		lastName := strings.TrimSpace(c.Request.FormValue("First_Name"))
		if firstName != "" && lastName != "" {
			dbContext.CreateMember(firstName, lastName)
		}

		http.Redirect(c.Writer, c.Request, "/", 301)
	})

	r.GET("update/:id", func(c *gin.Context) {
		memberId, _ := strconv.Atoi(c.Params.ByName("id"))
		member := dbContext.GetMember(memberId)
		obj := gin.H{"member": member}
		c.HTML(200, "update.tmpl", obj)
		// c.String(200, "["+myMember.Last_Name+"]")
	})

	r.POST("update/:id", func(c *gin.Context) {
		memberId, _ := strconv.Atoi(c.Params.ByName("id"))
		dbContext.UpdateMember(memberId, c.Request.FormValue("First_Name"), c.Request.FormValue("Last_Name"))
		http.Redirect(c.Writer, c.Request, "/", 301)
	})

	r.GET("delete/:id", func(c *gin.Context) {
		memberId, _ := strconv.Atoi(c.Params.ByName("id"))
		dbContext.DeleteMember(memberId)
		http.Redirect(c.Writer, c.Request, "/", 301)
	})

	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}
