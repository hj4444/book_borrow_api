package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type BaseModel struct {
	Id    uint      `json:"id" gorm:"primary_key"`
	CTime time.Time `json:"ctime" gorm:"column:ctime"`
	MTime time.Time `json:"mtime" gorm:"column:mtime"`
}

type Book struct {
	BaseModel
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Status      int
}
type Page struct {
	PageNum    int
	TotalPage  int
	TotalCount int64
}

type ResultPager struct {
	PageInfo *Page
	Datas    interface{} `json:"datas"`
}

type User struct {
	BaseModel
	Name       string `json:"name"`
	Password   string `json:"password"`
	ActualName string `json:"actualName" gorm:"column:actual_name"`
}
type UserBook struct {
	BaseModel
	BookId    uint   `json:"bookId" gorm:"column:book_id"`
	UserId    uint   `json:"userId" gorm:"column:user_id"`
	BeginTime string `json:"beginDate" gorm:"column:borrow_begin_date"`
	EndTime   string `json:"endDate" gorm:"column:borrow_end_date"`
}
type UserBookHistory struct {
	BaseModel
	BookName  string `json:"bookName" gorm:"column:name"`
	BeginTime string `json:"beginDate" gorm:"column:borrow_begin_date"`
	EndTime   string `json:"endDate" gorm:"column:borrow_end_date"`
	Url       string `json:"url"`
}
type ApiResult struct {
	Result    interface{} `json:"data,omitempty"`
	ErrorCode int         `json:"errorCode"`
	ErrorMsg  string      `json:"errorMsg"`
}

func DbSetup() {
	DBMS := "mysql"
	USER := "pointsystem"
	PASS := "qk365.com"
	PROTOCOL := "tcp(192.168.1.84:3306)"
	DBNAME := "books_borrow"
	var err error
	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=True&loc=Asia%2FTokyo"
	db, err = gorm.Open(DBMS, CONNECT)

	if err != nil {
		panic(err.Error())
	}
	db.LogMode(true)
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
}
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now()
		if createTimeField, ok := scope.FieldByName("CTime"); ok {
			if createTimeField.IsBlank {
				err := createTimeField.Set(nowTime)
				if err != nil {
					log.Fatal("err callback id ", err)
				}
			}
		}

		if modifyTimeField, ok := scope.FieldByName("MTime"); ok {
			if modifyTimeField.IsBlank {
				err := modifyTimeField.Set(nowTime)
				if err != nil {
					log.Fatal("err callback id ", err)
				}
			}
		}
	}
}

func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		err := scope.SetColumn("MTime", time.Now())
		if err != nil {
			log.Fatal("err callback id ", err)
		}
	}
}
func GetAllBooks(c echo.Context) error {
	var data []Book
	db.Debug().Where("status = 0").Find(&data)
	return c.JSON(http.StatusOK, ApiResult{Result: data, ErrorCode: 0, ErrorMsg: ""})
}

func Register(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return err
	}
	db.Debug().Create(&user)
	return c.JSON(http.StatusOK, ApiResult{Result: user, ErrorCode: 0, ErrorMsg: ""})
}

func Login(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return err
	}
	var newUser User
	db.Debug().Where("name=? and password=?", user.Name, user.Password).First(&newUser)
	if newUser.Id == 0 {
		return c.JSON(http.StatusOK, ApiResult{Result: newUser, ErrorCode: 1, ErrorMsg: "用户不存在"})
	}
	return c.JSON(http.StatusOK, ApiResult{Result: newUser, ErrorCode: 0, ErrorMsg: ""})
}

func BorrowBooks(c echo.Context) error {
	var userBooks []UserBook
	if err := c.Bind(&userBooks); err != nil {
		return err
	}
	tx := db.Begin()
	for i := 0; i < len(userBooks); i++ {
		var userBook = userBooks[i]
		tx.Debug().Create(&userBook)
		book := &Book{BaseModel: BaseModel{Id: userBook.BookId}}
		tx.Debug().Model(&book).UpdateColumn("status", 1).Where("id=?", userBook.BookId)
	}
	tx.Commit()
	return c.JSON(http.StatusOK, ApiResult{Result: "", ErrorCode: 0, ErrorMsg: ""})
}

func GetUserBorrowHistory(c echo.Context) error {
	id := c.Param("id")
	var data []UserBookHistory
	db.Debug().Table("user_book_map m").Joins("join book b on b.id=m.book_id").Joins("join user u on u.id=m.user_id").Select("b.name name,m.borrow_begin_date,m.borrow_end_date,b.url,m.ctime").Where("m.user_id=?", id).Scan(&data)
	return c.JSON(http.StatusOK, ApiResult{Result: data, ErrorCode: 0, ErrorMsg: ""})
}

func GetUserById(c echo.Context) error {
	id := c.Param("id")
	var user User
	db.Debug().Table("user").Select("id,name,password").Where("id=?", id).Scan(&user)
	return c.JSON(http.StatusOK, ApiResult{Result: user, ErrorCode: 0, ErrorMsg: ""})
}

func GetUserByName(c echo.Context) error {
	name := c.Param("name")
	var user User
	db.Debug().Table("user").Select("id,name,password").Where("name=?", name).Scan(&user)
	return c.JSON(http.StatusOK, ApiResult{Result: user, ErrorCode: 0, ErrorMsg: ""})
}

func (Book) TableName() string {
	return "book"
}

func (User) TableName() string {
	return "user"
}

func (UserBook) TableName() string {
	return "user_book_map"
}

var db *gorm.DB

func main() {
	DbSetup()
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/books", GetAllBooks)
	e.GET("/user/:id", GetUserById)
	e.GET("/user/:name", GetUserByName)
	e.POST("/user/login", Login)
	e.POST("/user/register", Register)
	e.POST("/user/borrow", BorrowBooks)
	e.GET("/user/borrowHistory/:id", GetUserBorrowHistory)
	e.Logger.Fatal(e.Start(":1323"))
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Print("shutdown server ....")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown:", err)
	}
	log.Print("server existing")
}
