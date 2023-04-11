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
			p := s.Peers[client]
			fmt.Printf("\n Peer %s disconnected", p.Name)
			delete(s.Peers, client)

		}
	}
}

func (s *Server) EventHandler(ev Event) {
	//add to event buffer
	//use mutex
	s.UiReadChan <- ev
	s.WriteEvent(ev)
}

// write event to client
func (s *Server) WriteEvent(ev Event) {
	b, err := json.Marshal(ev)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range s.Peers {
		if v.Name != ev.Body.From {
			_, err = v.Conn.Write(b)
			if err != nil {
				log.Fatal(err)
				continue
			}
		}
	}
}

func (s *Server) Accept() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Printf("\n Connection: %s", conn.RemoteAddr())
		peer := NewPeer(conn)
		s.Peers[conn.RemoteAddr().String()] = peer

		go s.HandleConn(conn)
	}
}

// read and listen to conn
func (s *Server) HandleConn(conn net.Conn) {
	defer func() {
		s.LeaveChan <- conn.RemoteAddr().String()
		conn.Close()
	}()
	buff := make([]byte, 2048)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Printf("Read error: %s", err.Error())
			break
		}

		msg := buff[:n]
		//fmt.Println(string(msg))
		var ev Event
		err = json.Unmarshal(msg, &ev)
		if err != nil {
			fmt.Printf("Unmarshal error: %s", err.Error())
		}

		//create event
		s.ReadChan <- ev
	}

}

func (s *Server) CreateEvent(m string, scope string) Event {
	msg := Message{From: s.Config.Host, Body: m}
	ev := Event{Scope: scope, Body: msg}
	return ev
}

func NewPeer(conn net.Conn) *Peer {
	return &Peer{Conn: conn}
}
