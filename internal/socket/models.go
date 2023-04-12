package socket

import (
	"net"
	"sync"

	"github.com/quinton11/chatline/internal/utils"
)

type Socket interface {
	GetRoomName() string
	GetWriteChan() chan Event
	GetUiReadChan() chan Event
	GetUserName() string
	CreateEvent(string, string) Event
}

type Server struct {
	Mu         sync.Mutex
	Config     utils.Room
	Listener   net.Listener
	Name       string
	Addr       string
	Peers      map[string]*Peer
	UiReadChan chan Event
	CloseChan  chan struct{}
	ReadChan   chan Event
	WriteChan  chan Event
	LeaveChan  chan string
}

type Client struct {
	User       string
	Room       utils.Room
	Conn       net.Conn
	Port       int
	UiReadChan chan Event
	ReadChan   chan Event
	WriteChan  chan Event
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
