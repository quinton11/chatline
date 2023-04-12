package cui

import (
	"fmt"
	"log"

	"github.com/fatih/color"

	"github.com/jroimartin/gocui"
	"github.com/quinton11/chatline/internal/socket"
)

var (
	blue    = color.New(color.FgBlue)
	green   = color.New(color.FgGreen)
	magenta = "\u001b[35m"
	red     = "\u001b[31m"
)

type Chatui struct {
	gui    *gocui.Gui
	socket socket.Socket
}

func NewChatUi(typ socket.Socket) *Chatui {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	return &Chatui{
		gui:    g,
		socket: typ,
	}
}

func (c *Chatui) Init() error {

	defer c.gui.Close()

	//set manager functions
	c.gui.SetManagerFunc(c.layout)

	//set keybindings
	if err := c.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, c.quit); err != nil {
		return err
	}

	if err := c.gui.SetKeybinding("chat", gocui.KeyEnter, gocui.ModNone, c.sendChat); err != nil {
		return err
	}

	go c.UpdateChats()

	//start mainloops

	if err := c.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}

func (c *Chatui) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	g.Cursor = true
	g.Mouse = true

	//chatline main frame
	if cline, err := g.SetView("chatline", 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		//config
		cline.Title = "Chatline"
		cline.Wrap = true
		cline.Autoscroll = true

	}

	//dim relative to chatline
	peerMaxX := int(float64(maxX) * 0.2)

	//Peers frame
	if p, err := g.SetView("peers", 1, 1, peerMaxX-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		//config
		p.Title = "Peers"
		p.Autoscroll = true
		p.Wrap = true

		rm := c.socket.GetRoomName()
		fmt.Fprintf(p, "Room: %s\n", rm)

	}

	//chat frame
	if ch, err := g.SetView("chats", peerMaxX, 1, maxX-2, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		//config
		ch.Title = "Chats"
		ch.Wrap = true
		ch.Autoscroll = true
		ch.Frame = true
	}

	//input msg frame
	if in, err := g.SetView("chat", peerMaxX, maxY-4, maxX-2, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		in.Title = "send"
		in.Wrap = true
		in.Autoscroll = true
		in.Editable = true
	}

	g.SetCurrentView("chat")

	return nil
}

func (c *Chatui) sendChat(g *gocui.Gui, v *gocui.View) error {
	if len(v.Buffer()) == 0 {
		v.SetCursor(0, 0)
		v.Clear()
		return nil
	}

	ev := c.socket.CreateEvent(v.Buffer(), socket.ChatEvent)
	c.socket.GetWriteChan() <- ev

	v.SetCursor(0, 0)
	v.Clear()
	return nil
}

func (c *Chatui) UpdateChats() {
	for {
		ev := <-c.socket.GetUiReadChan()
		c.gui.Update(func(g *gocui.Gui) error {
			view, _ := c.gui.View("chats")

			msg := formatEvent(ev, c.socket.GetUserName())
			fmt.Fprintln(view, msg)

			//check if event is init
			//if so update peers block
			//update with format room name \n Connected: number
			//connected increments on every init event and decreases on every leave event
			return nil
		})

	}
}

func formatEvent(ev socket.Event, you string) string {
	var msg string
	from := ev.Body.From
	switch {
	case ev.Scope == socket.InitEvent:
		if from == you {
			msg = fmt.Sprintf("%s[You] %s", magenta, ev.Body.Body)

		} else {
			msg = fmt.Sprintf("%s[%s] %s", magenta, from, ev.Body.Body)

		}
	case ev.Scope == socket.ChatEvent:
		if from == you {
			msg = blue.Sprintf("[You]: %s", ev.Body.Body)
		} else {
			msg = green.Sprintf("[%s]: %s", from, ev.Body.Body)
		}

	case ev.Scope == socket.LeaveEvent:
		msg = fmt.Sprintf("%s[%s] %s", red, from, ev.Body.Body)
	}
	return msg
}

func (c *Chatui) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
