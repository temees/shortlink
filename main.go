package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/teris-io/shortid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Url struct {
	ID   uint   `gorm:"primaryKey"`
	Src  string `json:"src" xml:"src" form:"src" query:"src"`
	Dest string `json:"dest" xml:"dest" form:"dest" query:"dest"`
}

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

var db *gorm.DB

func run() error {
	var err error
	db, err = gorm.Open(sqlite.Open("data.sqlite"), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&Url{})
	if err != nil {
		return err
	}

	e := echo.New()
	e.POST("/create", create)
	e.GET("/:id", redirect)

	err = e.Start(":3000")
	if err != nil {
		return err
	}
	return nil
}

func create(c echo.Context) error {
	url := new(Url)
	err := c.Bind(url)
	if err != nil {
		return err
	}
	url.Src, err = shortid.Generate()
	if err != nil {
		return err
	}
	db.Create(&url)
	log.Println(url)
	return c.JSON(http.StatusOK, url)
}

func redirect(c echo.Context) error {
	src := c.Param("id")
	url := &Url{Src: src}
	db.Where("src = ?", src).First(&url)
	log.Println(url.Dest)
	return c.Redirect(http.StatusMovedPermanently, url.Dest)
}
