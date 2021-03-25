package message

import (
	"github.com/diamondburned/arikawa/discord"
	"github.com/TheBlueOompaLoompa/gtkcord3/gtkcord/components/window"
	"github.com/TheBlueOompaLoompa/gtkcord3/gtkcord/gtkutils"
	"github.com/TheBlueOompaLoompa/gtkcord3/internal/log"
	"github.com/gotk3/gotk3/gtk"
	"github.com/atotto/clipboard"
	"strings"
)

func (m *Messages) menuAddAdmin(msg *Message, menu gtkutils.Container) {
	var canDelete = msg.AuthorID == m.c.Ready.User.ID
	if !canDelete {
		p, err := m.c.Permissions(m.GetChannelID(), m.c.Ready.User.ID)
		if err != nil {
			log.Errorln("Failed to get permissions on menuAddAdmin:", err)
			return
		}

		canDelete = p.Has(discord.PermissionManageMessages)
	}

	if canDelete {
		iDel, _ := gtk.MenuItemNewWithLabel("Delete Message")
		iDel.Connect("activate", func() {
			go func() {
				if err := m.c.DeleteMessage(m.GetChannelID(), msg.ID); err != nil {
					log.Errorln("Error deleting message:", err)
				}
			}()
		})
		iDel.Show()
		menu.Add(iDel)
	}

	if msg.AuthorID == m.c.Ready.User.ID {
		iEdit, _ := gtk.MenuItemNewWithLabel("Edit Message")
		iEdit.Connect("activate", func() {
			go func() {
				if err := m.Input.editMessage(msg.ID); err != nil {
					log.Errorln("Error editing message:", err)
				}
			}()
		})
		iEdit.Show()
		menu.Add(iEdit)
	}
}

func (m *Messages) menuAddReaction(msg *Message, menu gtkutils.Container) {
	canAddReactions := false
	isDM := false

	p, err := m.c.Permissions(m.GetChannelID(), m.c.Ready.User.ID)
	if err != nil {
		if(strings.Contains(err.Error(), "failed to get guild")){
			isDM = true
		}else{
			log.Errorln("Error getting permissions", err)
			return
		}
	}

	if !isDM {
		canAddReactions = p.Has(discord.PermissionAddReactions)
	}

	if canAddReactions || isDM {
		iReact, _ := gtk.MenuItemNewWithLabel("Add Reaction")
		iReact.Connect("activate", func() {
			go func() {
				reactionString, _ := clipboard.ReadAll()
				if err := m.c.State.React(m.GetChannelID(), msg.ID, reactionString); err != nil {
					log.Errorln("Error reacting to message:", err)
				}
			}()
		})
		iReact.Show()
		menu.Add(iReact)
	}
}

func (m *Messages) menuAddDebug(msg *Message, menu gtkutils.Container) {
	cpmsgID, _ := gtk.MenuItemNewWithLabel("Copy Message ID")
	cpmsgID.Connect("activate", func() {
		window.Window.Clipboard.SetText(msg.ID.String())
	})
	cpmsgID.Show()
	menu.Add(cpmsgID)

	cpchID, _ := gtk.MenuItemNewWithLabel("Copy Channel ID")
	cpchID.Connect("activate", func() {
		window.Window.Clipboard.SetText(m.GetChannelID().String())
	})
	cpchID.Show()
	menu.Add(cpchID)

	if m.GetGuildID().IsValid() {
		cpgID, _ := gtk.MenuItemNewWithLabel("Copy Guild ID")
		cpgID.Connect("activate", func() {
			window.Window.Clipboard.SetText(m.GetGuildID().String())
		})
		cpgID.Show()
		menu.Add(cpgID)
	}
}
