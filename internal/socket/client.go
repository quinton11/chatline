package socket

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/quinton11/chatline/internal/utils"
)

func NewClient(room utils.Room, name string) *Client {
	return &Client{
		Room:       room,
		Port:       5050,
		User:       name,
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

	//write init event with user name
	c.WriteChan <- c.CreateEvent("Joined Conversation", InitEvent)

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
			fmt.Println("Error in reading")
			fmt.Println(err)
			return err
		}

		msg := buff[:n]
		var ev Event
		err = json.Unmarshal(msg, &ev)
		if err != nil {
			continue
		}

		//
		c.ReadChan <- ev
	}
}

// on read write events
func (c *Client) Worker() {
	fmt.Println("Starting Worker")
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

func (c *Client) HandleRead(ev Event) {
	//push event to ui
	c.UiReadChan <- ev
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
	c.UiReadChan <- ev
	return nil
}

func (c *Client) CreateEvent(txt string, event string) Event {
	msg := Message{From: c.User, Body: txt}
	ev := Event{Scope: event, Body: msg}

	return ev
}

func (c *Client) GetRoomName() string {
	return c.Room.Name
}

// get write chan
func (c *Client) GetWriteChan() chan Event {
	return c.WriteChan
}

// get uiread chan
func (c *Client) GetUiReadChan() chan Event {
	return c.UiReadChan
}

// get user name
func (c *Client) GetUserName() string {
	return c.User
}
