package socket

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/quinton11/chatline/internal/utils"
)

func NewClient(room utils.Room) *Client {
	return &Client{
		Room:       room,
		Port:       5050,
		UiReadChan: make(chan Event),
		ReadChan:   make(chan Event),
		WriteChan:  make(chan Event),
	}
}

func (c *Client) Connect() error {
	addr := fmt.Sprintf("%s:%d", c.Room.Host, c.Room.Port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	defer conn.Close()

	go c.Worker()

	c.Conn = conn

	err = c.Read()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Read() error {
	buff := make([]byte, 2048)
	for {
		n, err := c.Conn.Read(buff)
		if err != nil {
			return err
		}

		msg := buff[:n]
		var ev Event
		err = json.Unmarshal(msg, &ev)
		if err != nil {
			fmt.Println("Couldn't parse event")
			continue
		}

		//
		c.ReadChan <- ev
	}
}

// remove worker
// call read write functionality in console ui
func (c *Client) Worker() {
	for {
		select {
		case read := <-c.ReadChan:
			c.HandleRead(read)

		case write := <-c.WriteChan:
			//send to server and print on command line
			err := c.HandleWrite(write)
			if err != nil {
				//Handle Error logging
				fmt.Println(err)
			}
		}

	}
}

func (c *Client) HandleRead(ev Event) error {
	//handle event
	//push mssage to buffer
	//use mutex
	c.UiReadChan <- ev

	return nil
}
func (c *Client) HandleWrite(ev Event) error {
	b, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	_, err = c.Conn.Write(b)
	if err != nil {
		return err
	}
	//add to buffer
	return nil
}

func (c *Client) CreateEvent(txt string, event string) Event {
	msg := Message{From: c.Conn.LocalAddr().String(), Body: txt}
	ev := Event{Scope: event, Body: msg}

	return ev
}
