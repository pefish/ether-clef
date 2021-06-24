package global

import "sync"

var (
	Password = ""
	AllowedMethod sync.Map
	Addresses sync.Map
)
