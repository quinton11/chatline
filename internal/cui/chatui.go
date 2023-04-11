package cui

import (
	"fmt"
	"log"

	"github.com/fatih/color"

	"github.com/jroimartin/gocui"
	"github.com/quinton11/chatline/internal/socket"
)

var (
	blue  = color.New(color.FgBlue)
	green = color.New(color.FgGreen)
)

type Chatui struct {
	gui      *gocui.Gui
	client   *socket.Client
	server   *socket.Server
	isserver bool
}

func NewChatUi(typ interface{}, isserver bool) *Chatui {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	if isserver {
		serv := typ.(*socket.Server)
		return &Chatui{
			gui:      g,
			server:   serv,
			isserver: isserver,
		}
	}

	cl := typ.(*socket.Client)
	return &Chatui{
		gui:      g,
		client:   cl,
		isserver: isserver,
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

	/* if !c.isserver {
		go c.UpdateChats()

	} */
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
		rm := ""
		if c.isserver {
			rm = c.server.Config.Name
		} else {
			rm = c.client.Room.Name
		}
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
	var ev socket.Event
	if c.isserver {
		ev = c.server.CreateEvent(v.Buffer(), socket.ChatEvent)
		//fmt.Println(ev)
		c.server.WriteChan <- ev
	} else {
		ev = c.client.CreateEvent(v.Buffer(), socket.ChatEvent)
		//fmt.Println(ev)
		c.client.WriteChan <- ev
	}

	v.SetCursor(0, 0)
	v.Clear()
	return nil
}

func (c *Chatui) UpdateChats() {
	for {
		if !c.isserver {
			ev := <-c.client.UiReadChan
			//instead use separate event buffer
			//allows you push both incoming and outgoing events to ui
			c.gui.Update(func(g *gocui.Gui) error {
				//sync mutex lock
				//loop through event buffer
				view, _ := c.gui.View("chats")

				msg := formatEvent(ev)

				//view.Clear()
				fmt.Fprintln(view, msg)
				return nil
			})
		} else {
			ev := <-c.server.UiReadChan
			c.gui.Update(func(g *gocui.Gui) error {
				//sync mutex lock
				//loop through event buffer
				view, _ := c.gui.View("chats")

				msg := formatEvent(ev)

				//view.Clear()
				fmt.Fprintln(view, msg)
				return nil
			})
		}

	}
}

func formatEvent(ev socket.Event) string {
	var msg string
	from := ev.Body.From
	switch {
	case ev.Scope == socket.InitEvent:
		fmt.Println("")
	case ev.Scope == socket.ChatEvent:
		//get from and msg
		//format: [from]: msg(blue)

		//if from is you : msg(green)
		if from == "You" {
			msg = blue.Sprintf("[%s]: %s", from, ev.Body.Body)
		} else {
			msg = green.Sprintf("[%s]: %s", from, ev.Body.Body)
		}

	case ev.Scope == socket.LeaveEvent:
		fmt.Println("")
	}
	return msg
}

func (c *Chatui) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
