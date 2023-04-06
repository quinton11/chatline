package socket

import (
	"github.com/quinton11/chatline/internal/utils"
)

// write and read to connection
type Socket interface {
}

type Server struct {
	Config utils.RoomConfig `json:"config"`
}
type Client struct {
	Config utils.RoomConfig `json:"config"`
}
