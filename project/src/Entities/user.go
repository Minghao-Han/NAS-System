package Entities

type User struct {
	UserId   int    `json:"id" form:"id"`
	UserName string `json:"username" form:"username"`
	Password []byte `json:"password" form:"password"`
	Capacity uint64 `json:"capacity" form:"capacity"`
	Margin   uint64 `json:"margin" form:"margin"`
}
