package Entities

type User struct {
	UserId   int    `json:"id" form:"id"`
	UserName string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Capacity int    `json:"capacity" form:"capacity"`
	Margin   int    `json:"margin" form:"margin"`
}
