package guild

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/gateway"
	"github.com/TheBlueOompaLoompa/gtkcord3/gtkcord/cache"
	"github.com/TheBlueOompaLoompa/gtkcord3/gtkcord/gtkutils"
	"github.com/TheBlueOompaLoompa/gtkcord3/gtkcord/ningen"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pkg/errors"
)

type GuildFolder struct {
	// Row that belongs to the parent list.
	*RevealerRow

	Icon *GuildFolderIcon
	Name string

	// Child list.
	List *gtk.ListBox

	Guilds   []*Guild
	Revealed bool
}

func newGuildFolder(
	s *ningen.State,
	folder gateway.GuildFolder,
	onSelect func(gf *GuildFolder, g *Guild)) (*GuildFolder, error) {

	if folder.Color == 0 {
		folder.Color = 0x7289DA
	}

	guildList, _ := gtk.ListBoxNew()
	guildList.SetActivateOnSingleClick(true)

	var Folder = &GuildFolder{
		List:   guildList,
		Name:   folder.Name,
		Guilds: make([]*Guild, 0, len(folder.GuildIDs)),
	}

	// Bind the child list independent of the parent list.
	guildList.Connect("row-activated", func(l *gtk.ListBox, r *gtk.ListBoxRow) {
		i := r.GetIndex()
		Folder.unselectAll(i)

		row := Folder.Guilds[i]
		row.Unread.SetActive(true)
		onSelect(Folder, row)
	})

	// Used to mark read and unread.
	var unread, pinged bool

	for _, id := range folder.GuildIDs {
		r, err := newGuildRow(s, id, nil, Folder)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to load guild "+id.String())
		}
		r.Parent = Folder
		Folder.Guilds = append(Folder.Guilds, r)

		Folder.List.Add(r)
	}

	// Take care of the icon part.
	icon := newGuildFolderIcon(Folder.Guilds)
	Folder.Icon = icon

	// On click, toggle revealer.
	rev := newRevealerRow(icon, guildList, func(reveal bool) {
		// Expand/collapse the icon
		icon.setReveal(reveal)
	})

	Folder.RevealerRow = rev
	gtkutils.InjectCSSUnsafe(rev, "guild-folder", "")
	// Folder.Style, _ = rev.GetStyleContext()
	// Folder.Style.AddClass("guild-folder")

	// Show name on hover.
	BindNameDirect(rev.Button, rev.Strip, &Folder.Name)

	// Set the unread status for parent.
	switch {
	case pinged:
		rev.Strip.SetPinged()
	case unread:
		rev.Strip.SetUnread()
	}

	// Color time.
	color := fmt.Sprintf("#%06X", folder.Color.Uint32())

	// Color the folder icon.
	gtkutils.InjectCSSUnsafe(icon.Folder, "", `* { color: `+color+`; }`)
	// Color the collapsed folder background.
	gtkutils.AddCSSUnsafe(icon.Style, `
		*.collapsed {
			/* We have to use mix because alpha breaks with border-radius */
			background-color: mix(@theme_bg_color, `+color+`, 0.4);
		}
	`)

	// Add some room:
	rev.ListBoxRow.SetSizeRequest(IconSize+IconPadding*3, IconSize+IconPadding/2)
	gtkutils.Margin2(rev, IconPadding/2, 0)

	return Folder, nil
}

func (f *GuildFolder) unselectAll(except int) {
	for i, r := range f.Guilds {
		if i == except {
			continue
		}
		r.Unread.SetActive(false)
	}
}

func (f *GuildFolder) setUnread(unread, pinged bool) {
	// Check all children guilds
	if !unread || !pinged {
		for _, g := range f.Guilds {
			switch _unread, _pinged := g.Unread.State(); {
			case _pinged:
				pinged = true
				fallthrough
			case _unread:
				unread = true
			}
		}
	}

	switch {
	case pinged:
		f.Strip.SetPinged()
	case unread:
		f.Strip.SetUnread()
	default:
		f.Strip.SetRead()
	}
}

type GuildFolderIcon struct {
	folder []*Guild

	// Main stack, switches between "guilds" and "folder"
	*gtk.Stack
	Style *gtk.StyleContext

	Guilds *gtk.Grid     // contains 4 images always.
	Images [4]*gtk.Image // first 4 of folder.Guilds

	Folder *gtk.Image
}

func newGuildFolderIcon(guilds []*Guild) *GuildFolderIcon {
	icon := &GuildFolderIcon{
		folder: guilds,
	}

	icon.Stack, _ = gtk.StackNew()
	icon.Stack.SetTransitionType(gtk.STACK_TRANSITION_TYPE_SLIDE_UP) // unsure
	icon.Stack.SetSizeRequest(IconSize, IconSize)

	icon.Style, _ = icon.Stack.GetStyleContext()
	icon.Style.AddClass("collapsed") // used for coloring

	icon.Folder, _ = gtk.ImageNew()
	gtkutils.ImageSetIcon(icon.Folder, "folder-symbolic", FolderSize)

	icon.Guilds, _ = gtk.GridNew()
	icon.Guilds.SetHAlign(gtk.ALIGN_CENTER)
	icon.Guilds.SetVAlign(gtk.ALIGN_CENTER)
	icon.Guilds.SetRowSpacing(4) // calculated from Discord
	icon.Guilds.SetRowHomogeneous(true)
	icon.Guilds.SetColumnSpacing(4)
	icon.Guilds.SetColumnHomogeneous(true)

	// Make dummy images.
	for i := range icon.Images {
		img, _ := gtk.ImageNew()
		img.SetSizeRequest(16, 16)

		icon.Images[i] = img
	}

	// Set the dummy images in a grid.
	// [0] [1]
	// [2] [3]
	icon.Guilds.Attach(icon.Images[0], 0, 0, 1, 1)
	icon.Guilds.Attach(icon.Images[1], 1, 0, 1, 1)
	icon.Guilds.Attach(icon.Images[2], 0, 1, 1, 1)
	icon.Guilds.Attach(icon.Images[3], 1, 1, 1, 1)

	// Asynchronously fetch the icons.
	for i := 0; i < len(guilds) && i < 4; i++ {
		url := guilds[i].IURL
		if url == "" {
			continue
		}
		// Replace GIF with PNG to save CPU cycles.
		url = strings.Replace(url, "gif", "png", -1)
		url += "?size=64" // same as guild.go

		cache.AsyncFetchUnsafe(url, icon.Images[i], 16, 16, cache.Round)
	}

	// Add things together.
	icon.Stack.AddNamed(icon.Guilds, "guilds")
	icon.Stack.AddNamed(icon.Folder, "folder")

	return icon
}

// called with revealer
func (i *GuildFolderIcon) setReveal(reveal bool) {
	if reveal {
		// show the folder.
		i.Stack.SetVisibleChildName("folder")
		i.Style.RemoveClass("collapsed")
	} else {
		// show the guilds
		i.Stack.SetVisibleChildName("guilds")
		i.Style.AddClass("collapsed")
	}
}

type RevealerRow struct {
	*gtk.ListBoxRow
	Strip    *UnreadStrip
	Button   *gtk.Button
	Revealer *gtk.Revealer
}

func newRevealerRow(button, reveal gtk.IWidget, click func(reveal bool)) *RevealerRow {
	r, _ := gtk.RevealerNew()
	r.Show()
	r.SetTransitionType(gtk.REVEALER_TRANSITION_TYPE_SLIDE_UP)
	r.SetRevealChild(false)
	r.Add(reveal)

	btn, _ := gtk.ButtonNew()
	btn.Show()
	btn.SetHAlign(gtk.ALIGN_CENTER)
	btn.SetVAlign(gtk.ALIGN_CENTER)
	btn.SetRelief(gtk.RELIEF_NONE)
	btn.Add(button)

	// Wrap both the widget child and the revealer
	b, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	b.Show()
	b.Add(btn)
	b.Add(r)

	// Wrap the stack inside the unread strip overlay.
	strip := NewUnreadStrip(b)

	row, _ := gtk.ListBoxRowNew()
	row.Show()
	row.SetActivatable(true)
	row.SetSelectable(false)
	row.SetHAlign(gtk.ALIGN_CENTER)
	row.SetVAlign(gtk.ALIGN_CENTER)
	row.Add(strip)

	btn.Connect("clicked", func() {
		reveal := !r.GetRevealChild()
		r.SetRevealChild(reveal)
		click(reveal)
		strip.SetSuppress(reveal)
	})

	return &RevealerRow{row, strip, btn, r}
}
