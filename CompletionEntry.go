package main

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CompletionEntry is an Entry with options displayed in a PopUpMenu.
type CompletionEntry struct {
	widget.Entry
	popupMenu     *widget.PopUp
	navigableList *navigableList
	Options       []string
	pause         bool
	itemHeight    float32
}

// NewCompletionEntry creates a new CompletionEntry which creates a popup menu that responds to keystrokes to navigate through the items without losing the editing ability of the text input.
func NewCompletionEntry(options []string) *CompletionEntry {
	c := &CompletionEntry{Options: options}
	c.ExtendBaseWidget(c)
	return c
}

// HideCompletion hides the completion menu.
func (c *CompletionEntry) HideCompletion() {
	if c.popupMenu != nil {
		c.popupMenu.Hide()
	}
}

// Move changes the relative position of the select entry.
//
// Implements: fyne.Widget
func (c *CompletionEntry) Move(pos fyne.Position) {
	c.Entry.Move(pos)
	if c.popupMenu != nil {
		c.popupMenu.Resize(c.maxSize())
		c.popupMenu.Move(c.popUpPos())
	}
}

// Refresh the list to update the options to display.
func (c *CompletionEntry) Refresh() {
	c.Entry.Refresh()
	if c.navigableList != nil {
		c.navigableList.SetOptions(c.Options)
	}
}

// SetOptions set the completion list with itemList and update the view.
func (c *CompletionEntry) SetOptions(itemList []string) {
	c.Options = itemList
	c.Refresh()
}

// ShowCompletion displays the completion menu
func (c *CompletionEntry) ShowCompletion() {
	if c.pause {
		return
	}
	if len(c.Options) == 0 {
		c.HideCompletion()
		return
	}

	if c.navigableList == nil {
		c.navigableList = newNavigableList(c.Options, &c.Entry, c.setTextFromMenu, c.HideCompletion)
	}
	holder := fyne.CurrentApp().Driver().CanvasForObject(c)

	if c.popupMenu == nil {
		c.popupMenu = widget.NewPopUp(c.navigableList, holder)
	}
	c.popupMenu.Resize(c.maxSize())
	c.popupMenu.ShowAtPosition(c.popUpPos())
	holder.Focus(c.navigableList)
}

// calculate the max size to make the popup to cover everything below the entry
func (c *CompletionEntry) maxSize() fyne.Size {
	cnv := fyne.CurrentApp().Driver().CanvasForObject(c)

	if c.itemHeight == 0 {
		// set item height to cache
		c.itemHeight = c.navigableList.CreateItem().MinSize().Height
	}

	listheight := float32(len(c.Options))*(c.itemHeight+2*theme.Padding()+theme.SeparatorThicknessSize()) + 2*theme.Padding()
	canvasSize := cnv.Size()
	entrySize := c.Size()
	if canvasSize.Height > listheight {
		return fyne.NewSize(entrySize.Width, listheight)
	}

	return fyne.NewSize(
		entrySize.Width,
		canvasSize.Height-c.Position().Y-entrySize.Height-theme.InputBorderSize()-theme.Padding())
}

// calculate where the popup should appear
func (c *CompletionEntry) popUpPos() fyne.Position {
	entryPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(c)
	return entryPos.Add(fyne.NewPos(0, c.Size().Height))
}

// Prevent the menu to open when the user validate value from the menu.
func (c *CompletionEntry) setTextFromMenu(s string) {
	c.pause = true
	c.Entry.SetText(s)
	c.Entry.CursorColumn = len([]rune(s))
	c.Entry.Refresh()
	c.pause = false
	c.popupMenu.Hide()
}

type navigableList struct {
	widget.List
	entry           *widget.Entry
	selected        int
	setTextFromMenu func(string)
	hide            func()
	navigating      bool
	items           []string
}

////Only changed function
func newNavigableList(items []string, entry *widget.Entry, setTextFromMenu func(string), hide func()) *navigableList {
	n := &navigableList{
		entry:           entry,
		selected:        -1,
		setTextFromMenu: setTextFromMenu,
		hide:            hide,
		items:           items,
	}

	n.List = widget.List{
		Length: func() int {
			return len(n.items)
		},
		CreateItem: func() fyne.CanvasObject {

			city := widget.NewRichTextFromMarkdown("City Lowercase")
			city.Move(fyne.NewPos(0, 0))

			location := canvas.NewText("COUNTRY - BST", color.NRGBA{0xFF, 0xFF, 0xFF, 0xBF})
			location.TextStyle.Monospace = true
			location.TextSize = 10
			location.Move(fyne.NewPos(theme.Padding()*2, -theme.Padding()*2))

			return container.NewVBox(city, container.NewWithoutLayout(location))

		},
		UpdateItem: func(i widget.ListItemID, o fyne.CanvasObject) {

			//bug catch
			if i > len(n.items)-1 {
				fmt.Println("Crashes if not caught here")
				return
			}

			//would be nice to pass city struct in here instead of splitting a string
			c := o.(*fyne.Container)
			split := strings.Split(n.items[i], "--")

			city := c.Objects[0].(*widget.RichText)
			city.ParseMarkdown(split[0])

			countryAndTZ := c.Objects[1].(*fyne.Container).Objects[0].(*canvas.Text)
			z, _ := time.LoadLocation(split[2])
			t := time.Now().In(z)
			countryAndTZ.Text = (strings.ToUpper(split[1]) + " · " + t.Format("MST"))

		},
		OnSelected: func(id widget.ListItemID) {
			if !n.navigating && id > -1 {
				entry.OnSubmitted(n.items[id])
			}
			n.navigating = false
		},
	}
	n.ExtendBaseWidget(n)
	return n
}

// Implements: fyne.Focusable
func (n *navigableList) FocusGained() {
}

// Implements: fyne.Focusable
func (n *navigableList) FocusLost() {
}

func (n *navigableList) SetOptions(items []string) {
	n.Unselect(n.selected)
	n.items = items
	n.Refresh()
	n.selected = -1
}

func (n *navigableList) TypedKey(event *fyne.KeyEvent) {
	switch event.Name {
	case fyne.KeyDown:
		if n.selected < len(n.items)-1 {
			n.selected++
		} else {
			n.selected = 0
		}
		n.navigating = true
		n.Select(n.selected)

	case fyne.KeyUp:
		if n.selected > 0 {
			n.selected--
		} else {
			n.selected = len(n.items) - 1
		}
		n.navigating = true
		n.Select(n.selected)
	case fyne.KeyReturn, fyne.KeyEnter:
		if n.selected == -1 { // so the user want to submit the entry
			n.hide()
			n.entry.TypedKey(event)
		} else {
			n.navigating = false
			n.OnSelected(n.selected)
		}
	case fyne.KeyEscape:
		n.hide()
	default:
		n.entry.TypedKey(event)

	}
}

func (n *navigableList) TypedRune(r rune) {
	n.entry.TypedRune(r)
}