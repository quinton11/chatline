package socket

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/quinton11/chatline/internal/utils"
)

// write and read to connection

func NewServer(r utils.Room) *Server {
	addr := fmt.Sprintf(":%d", r.Port)
	return &Server{
		Config:     r,
		Addr:       addr,
		Name:       "Host",
		Peers:      make(map[string]*Peer),
		UiReadChan: make(chan Event),
		CloseChan:  make(chan struct{}),
		ReadChan:   make(chan Event),
		WriteChan:  make(chan Event),
		LeaveChan:  make(chan string),
	}
}

// start tcp server
func (s *Server) Init() error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	defer listener.Close()

	s.Listener = listener

	//orchestrate - set up channels to listen
	go s.Orchestra()

	//accept connections
	go s.Accept()

	<-s.CloseChan

	return nil
}

// set up chans
func (s *Server) Orchestra() {
	for {
		select {
		case read := <-s.ReadChan:
			s.EventHandler(read)
		case write := <-s.WriteChan:
			//pass to ui
			s.UiReadChan <- write
			s.WriteEvent(write)

		case client := <-s.LeaveChan:
			//remove client from pool
			//p := s.Peers[client]
			s.Mu.Lock()
			delete(s.Peers, client)
			s.Mu.Unlock()

		}
	}
}

func (s *Server) EventHandler(ev Event) {
	//if scope is init event, update peer name
	s.UiReadChan <- ev
	s.WriteEvent(ev)
}

// write event to client
func (s *Server) WriteEvent(ev Event) {
	b, err := json.Marshal(ev)
	if err != nil {
		log.Fatal(err)
	}
	s.Mu.Lock()
	for _, v := range s.Peers {
		if v.Name != ev.Body.From {
			_, err = v.Conn.Write(b)
			if err != nil {
				log.Fatal(err)
				continue
			}
		}
	}
	s.Mu.Unlock()
}

func (s *Server) Accept() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Println(err)
		}

		peer := NewPeer(conn)

		s.Mu.Lock()
		s.Peers[conn.RemoteAddr().String()] = peer
		s.Mu.Unlock()

		go s.HandleConn(conn)
	}
}

// read and listen to conn
func (s *Server) HandleConn(conn net.Conn) {
	defer func() {
		conn.Close()
	}()
	buff := make([]byte, 2048)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			break
		}

		msg := buff[:n]

		var ev Event
		err = json.Unmarshal(msg, &ev)
		if err != nil {
			fmt.Printf("Unmarshal error: %s", err.Error())
		}

		if ev.Scope == InitEvent {
			//update peer name
			s.Mu.Lock()
			peer := s.Peers[conn.RemoteAddr().String()]
			peer.Name = ev.Body.From
			s.Mu.Unlock()
		}

		//create event
		s.ReadChan <- ev
	}
	//create leave event, push to relevant channels
	s.LeaveChan <- conn.RemoteAddr().String()
	p := s.Peers[conn.RemoteAddr().String()]
	s.WriteChan <- s.LeaveEvent(p.Name)

}

func (s *Server) InitEvent(client string) Event {
	msg := Message{From: client, Body: "Joined Conversation"}
	ev := Event{Scope: InitEvent, Body: msg}

	return ev
}

func (s *Server) LeaveEvent(client string) Event {
	msg := Message{From: client, Body: "Left Conversation"}
	ev := Event{Scope: LeaveEvent, Body: msg}

	return ev
}

func (s *Server) CreateEvent(m string, scope string) Event {
	msg := Message{From: s.Name, Body: m}
	ev := Event{Scope: scope, Body: msg}
	return ev
}

func (s *Server) GetRoomName() string {
	return s.Config.Name
}

// get write chan
func (s *Server) GetWriteChan() chan Event {
	return s.WriteChan
}

// get uiread chan
func (s *Server) GetUiReadChan() chan Event {
	return s.UiReadChan
}

// get user name
func (s *Server) GetUserName() string {
	return s.Name
}

func NewPeer(conn net.Conn) *Peer {
	return &Peer{Conn: conn}
}
