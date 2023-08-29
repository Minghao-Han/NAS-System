package fsp

type Port struct {
	csPort      int
	dsPort      int
	connections []Connection
}
