package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"nas/app/utils"
	"net/http"
	"strconv"
)

func main() {
	//创建一个路由Handler
	//router := gin.Default()
	user := utils.User{
		UserId:   0,
		UserName: "aaaa",
		Password: "bbbb",
		Capacity: 0,
		Margin:   0,
	}
	//插入数据测试
	//_, err := utils.Insert(user)
	//if err != nil {
	//	return
	//}

	////查询数据测试
	//users, err := utils.Query()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("查询数据")
	//fmt.Println(users)
	//user.Capacity = 9
	////修改数据测试
	//_, err = utils.Update(user)
	//user2, err := utils.Query()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("修改数据")
	//fmt.Println(user2)
	////删除数据测试
	//_, err = utils.Del(0)
	//user3, err := utils.Query()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("删除数据")
	//fmt.Println(user3)

	//创建路由器
	router := gin.Default()

	//get方法的查询
	router.GET("/query", func(c *gin.Context) {
		//utils.Query()

		c.JSON(http.StatusOK, gin.H{
			"result": user,
		})
	})

	//利用post方法新增数据
	router.POST("/add", func(c *gin.Context) {
		var u utils.User
		err := c.Bind(&u)
		if err != nil {
			log.Fatal(err)
		}
		Id, err := utils.Insert(u)
		fmt.Print("id=", Id)
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("%s 插入成功", u.UserName),
		})
	})

	//利用put方法修改数据
	router.PUT("/update", func(c *gin.Context) {
		var u utils.User
		err := c.Bind(&u)
		if err != nil {
			log.Fatal(err)
		}
		num, err := utils.Update(u)
		fmt.Print("num=", num)
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("修改id: %d 成功", u.UserId),
		})
	})

	//利用DELETE请求方法通过id删除
	router.DELETE("/delete/:id", func(c *gin.Context) {
		id := c.Param("id")

		Id, err := strconv.Atoi(id)

		if err != nil {
			log.Fatalln(err)
		}
		rows, err := utils.Del(Id)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("delete rows ", rows)

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully deleted user: %s", id),
		})
	})

	router.Run(":8000")
}
