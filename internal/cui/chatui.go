package cui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/quinton11/chatline/internal/socket"
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
		c.server.WriteChan <- ev
	} else {
		ev = c.client.CreateEvent(v.Buffer(), socket.ChatEvent)
		c.client.WriteChan <- ev
	}

	v.SetCursor(0, 0)
	v.Clear()
	return nil
}

func (c *Chatui) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
