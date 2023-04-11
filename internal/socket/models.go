package socket

import (
	"net"

	"github.com/quinton11/chatline/internal/utils"
)

type Socket interface {
}

type Server struct {
	Config    utils.Room
	Listener  net.Listener
	Addr      string
	Peers     map[string]*Peer
	EvBuff    []Event
	CloseChan chan struct{}
	ReadChan  chan Event
	WriteChan chan Event
	LeaveChan chan string
}

type Client struct {
	Room      utils.Room
	Conn      net.Conn
	Port      int
	EvBuff    []Event
	ReadChan  chan Event
	WriteChan chan Event
}

type Peer struct {
	Name string `json:"name"`
	Conn net.Conn
}

type Event struct {
	Scope string  `json:"scope"`
	Body  Message `json:"body"`
}

type Message struct {
	From string      `json:"from"`
	Body interface{} `json:"body"`
}

const (
	InitEvent  = "init"
	ChatEvent  = "chat"
	LeaveEvent = "leave"
)
