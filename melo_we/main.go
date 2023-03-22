package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
	"io"
	"net/http"
	"os"
)

type com struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

func IsFile_exist(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	} else {
		return !info.IsDir()
	}
}

func CreateFile(fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()
}
func WriteFile(fileName string, data com) {
	file, op_err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if op_err != nil {
		fmt.Println(op_err.Error())
		return
	}
	defer file.Close()

	existingData := []interface{}{}

	if fileInfo, err := file.Stat(); err == nil && fileInfo.Size() > 0 {
		// The file already has data, so we need to decode it into the existingData slice
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&existingData)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	existingData = append(existingData, data)

	// Seek back to the beginning of the file before writing the updated data
	file.Seek(0, io.SeekStart)
	encoder := json.NewEncoder(file)
	e_err := encoder.Encode(existingData)

	if e_err != nil {
		fmt.Println((e_err.Error()))
		return
	}
}
func ReadFile_byLine(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var x []com

	err = json.Unmarshal(data, &x)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, p := range x {
		username := p.Name
		text := p.Text

		fmt.Println(username, " ", text)
		//c.HTML(http.StatusOK, "web.html", gin.H{
		//	"Name": username,
		//	"Text": text,
		//})
	}

}

func main() {
	fileName := "cmt.json"

	//set server
	ginServer := gin.Default()
	ginServer.Use(favicon.New("./melo.ico"))

	//load static source for html
	ginServer.LoadHTMLFiles("./web/web.html", "./web/web_c.html")
	ginServer.Static("/features", "./features")
	ginServer.Static("cmt.json", ".")
	ginServer.Use(static.Serve("/", static.LocalFile("public", true)))
	//request response
	ginServer.GET("/melo", func(c *gin.Context) {
		c.HTML(http.StatusOK, "web.html", gin.H{})
	})

	ginServer.GET("/melo/cmt.json", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "cmt.json")
	})

	ginServer.POST("/melo", func(c *gin.Context) {
		username := c.PostForm("username")
		cot := c.PostForm("comment")
		cot_data := com{Name: username, Text: cot}
		WriteFile(fileName, cot_data)

		c.HTML(http.StatusOK, "web.html", gin.H{})

		//ReadFile_byLine(fileName, c)
	})

	//server address
	ginServer.Run(":1234")

}
