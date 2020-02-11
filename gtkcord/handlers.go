package gtkcord

import (
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
)

func (a *Application) hookEvents() {
	a.State.AddHandler(func(v interface{}) {
		a.busy.Lock()
		defer a.busy.Unlock()

		switch v := v.(type) {
		case *gateway.MessageCreateEvent:
			a.onMessageCreate(v)
		case *gateway.MessageUpdateEvent:
			a.onMessageUpdate(v)
		case *gateway.MessageDeleteEvent:
			a.onMessageDelete(v)
		case *gateway.MessageDeleteBulkEvent:
			a.onMessageDeleteBulk(v)
		}
	})
}

func (a *Application) onMessageCreate(m *gateway.MessageCreateEvent) {
	mw, ok := a.Messages.(*Messages)
	if !ok {
		return
	}

	if m.ChannelID != mw.ChannelID {
		return
	}

	if err := mw.Insert(a.State, a.parser, discord.Message(*m)); err != nil {
		logWrap(err, "Failed to insert message from "+m.Author.Username)
	}
}

func (a *Application) onMessageUpdate(m *gateway.MessageUpdateEvent) {
	mw, ok := a.Messages.(*Messages)
	if !ok {
		return
	}

	if m.ChannelID != mw.ChannelID {
		return
	}

	mw.Update(a.State, a.parser, discord.Message(*m))
}

func (a *Application) onMessageDelete(m *gateway.MessageDeleteEvent) {
	mw, ok := a.Messages.(*Messages)
	if !ok {
		return
	}

	if m.ChannelID != mw.ChannelID {
		return
	}

	mw.Delete(m.ID)
}

func (a *Application) onMessageDeleteBulk(m *gateway.MessageDeleteBulkEvent) {
	mw, ok := a.Messages.(*Messages)
	if !ok {
		return
	}

	if m.ChannelID != mw.ChannelID {
		return
	}

	for _, id := range m.IDs {
		mw.Delete(id)
	}
}
