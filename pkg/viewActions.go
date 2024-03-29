package pkg

/*
import (
	"errors"
	"fmt"

	"github.com/jroimartin/gocui"
	aur "github.com/logrusorgru/aurora"
)

var popupDisplayed = false

// togglePopup toggles the popup window's visibility by bringing it to
// either the very front or very back.
func togglePopup() {
	if !popupDisplayed {
		g.SetViewOnTop(popup)
		g.SetViewOnTop(searchOutline)
		g.SetViewOnTop(searchSymbol)
		g.SetViewOnTop(search)
		g.SetCurrentView(search)
		setTitle(searchOutline, "Enter brewery and beer name...")
	} else {
		setTitle(popup, "")
		g.SetViewOnBottom(popup)
		g.SetViewOnBottom(searchSymbol)
		g.SetViewOnBottom(searchOutline)
		g.SetViewOnBottom(search)
		g.SetCurrentView(input)
	}

	popupDisplayed = !popupDisplayed
}

// updatePromptSymbol highlights the selected mode.
func updatePromptSymbol() {
	v, _ := g.View(promptSymbol)
	v.Clear()
	switch mode := c.GetMode(); mode {
	case stocking:
		fmt.Fprintf(v, "%s >>", aur.BgBrown("Stocking"))
	case serving:
		fmt.Fprintf(v, "%s >>", aur.BgGreen("Serving"))
	}

}

// clearView clears a given gocui view and hides the cursor.
func clearView(view string) {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View(view)
		if err != nil {
			return err
		}
		v.Clear()
		x, y := v.Cursor()
		v.MoveCursor(-x, -y, true)

		return nil
	})
}

// promptEditor handles keyboard input for the input line of a prompt.
func promptEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	if ch != 0 && mod == 0 {
		v.EditWrite(ch)
		return
	}

	switch key {
	case gocui.KeySpace:
		v.EditWrite(' ')
	case gocui.KeyBackspace, gocui.KeyBackspace2:
		v.EditDelete(true)
	case gocui.KeyDelete:
		v.EditDelete(false)
	case gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case gocui.KeyArrowDown:
		_ = v.SetCursor(len(v.Buffer())-1, 0)
	case gocui.KeyArrowUp:
		v.MoveCursor(0, -1, false)
	case gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case gocui.KeyArrowRight:

	}
}

// displayError displays an error message in a bright red popup.
//
// TODO: Make it possible to escape from this popup (using hideError)
func displayError(e error) error {
	maxX, maxY := g.Size()
	x0 := maxX / 6
	y0 := maxY / 6
	x1 := (5 * maxX) / 6
	y1 := (5 * maxY) / 6

	if v, err := g.SetView(errorView, x0, y0, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "ERROR"
		v.Frame = true
		v.Wrap = true
		v.Autoscroll = true
		v.BgColor = gocui.ColorRed
		v.FgColor = gocui.ColorWhite

		v.Clear()
		fmt.Fprintln(v, e.Error())
		g.SetCurrentView(v.Name())
	}

	return nil
}

// hideError hides the error view.
//
// TODO: figure out why this isn't being called by the ESC keybinding.
func hideError(g *gocui.Gui, _ *gocui.View) error {
	g.DeleteView(errorView)
	return nil
}

// popupScrollUp moves the cursor selection a row up in the popup.
func popupScrollUp(_ *gocui.Gui, v *gocui.View) error {
	err := moveViewCursorUp(v)
	if err != nil {
		logAllError(err)
	}
	return nil
}

// popupScrollDown moves the cursor selection a row down in the popup.
func popupScrollDown(_ *gocui.Gui, v *gocui.View) error {
	err := moveViewCursorDown(v, false)
	if err != nil {
		logAllError(err)
	}
	return err
}

// moveViewCursorDown moves the cursor selection a row down in a given view.
func moveViewCursorDown(v *gocui.View, allowEmpty bool) error {
	cx, cy := v.Cursor()
	ox, oy := v.Origin()
	nextLine, err := getNextViewLine(v)
	if err != nil {
		return err
	}
	if !allowEmpty && nextLine == "" {
		return nil
	}
	if err := v.SetCursor(cx, cy+1); err != nil {
		if err := v.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}
	return nil
}

// moveViewCursorUp moves the cursor selection a row up in a given view.
func moveViewCursorUp(v *gocui.View) (err error) {
	cx, cy := v.Cursor()
	ox, oy := v.Origin()
	switch {
	case cy == 0 && oy == 0: // already at the top
		return nil
	case cy > 0: // cursor has priority over origin
		if err = v.SetCursor(cx, cy-1); err != nil {
			return err
		}
		return nil
	case oy > 0:
		if err = v.SetOrigin(ox, oy-1); err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid cursor or origin position")
	}
}

// getViewLine fetches the text shown in the row highlighted by
// the user's cursor position.
func getViewLine(v *gocui.View) (string, error) {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	return l, err
}

// getViewLine fetches the text shown in the row belwo the one
// highlighted by the user's cursor position.
func getNextViewLine(v *gocui.View) (string, error) {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy + 1); err != nil {
		l = ""
	}

	return l, err
}

// resetViewCursor sets the user's cursor position to the topmost
// row of a given gocui view.
func resetViewCursor(v *gocui.View) error {
	ox, _ := v.Origin()
	cx, _ := v.Cursor()
	if err := v.SetCursor(ox, 0); err != nil {
		if err := v.SetOrigin(cx, 0); err != nil {
			return err
		}
	}
	return nil
}

// scrollView scrolls the user's current highlighted selection in a gocui view
// by the given number of rows.
//
// Positive dy corresponds to a downward movement, while negative dy moves the
// cursor up. If the cursor cannot be moved, it is clipped to the maximum or
// minimum row.
func scrollView(v *gocui.View, dy int) error {
	if v != nil {
		v.Autoscroll = false
		ox, oy := v.Origin()
		_, height := v.Size()
		if dy > 0 {
			if l, _ := v.Line(height); l == "" {
				return nil
			}
		}
		if err := v.SetOrigin(ox, oy+dy); err != nil {
			return err
		}
	}
	return nil
}

// setTitle sets the title of a view in the main gocui object.
func setTitle(view string, title string) {
	v, _ := g.View(view)
	v.Title = title
}
*/
