package window

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

type Window interface {
	Conn() *xgb.Conn
	XRoot() xproto.Window
	XWindow() xproto.Window
	Close()
	Kill()
	SetXBorderWidth(uint32)
	Move(int16, int16)
	Resize(uint16, uint16)
	MoveResize(int16, int16, uint16, uint16)
	Raise()
	Lower()
	Stack(xproto.Window, uint32)
	Above(Window)
	Below(Window)
	Hide()
	Show()
}

func New(c *xgb.Conn, w xproto.Window, r xproto.Window) Window {
	return &window{c, w, r}
}

type window struct {
	c *xgb.Conn
	w xproto.Window
	r xproto.Window
}

func (w *window) Conn() *xgb.Conn {
	return w.c
}

func (w *window) XRoot() xproto.Window {
	return w.r
}

func (w *window) XWindow() xproto.Window {
	return w.w
}

func (w *window) Close() {
	//send_client_message(w.Window, ewmh->WM_PROTOCOLS, WM_DELETE_WINDOW);
}

func (w *window) Kill() {
	xproto.KillClientChecked(w.c, uint32(w.w))
}

func (w *window) SetXBorderWidth(bw uint32) {
	xproto.ConfigureWindowChecked(w.c, w.w, xproto.ConfigWindowBorderWidth, []uint32{bw})
}

func (w *window) Move(x, y int16) {
	xproto.ConfigureWindowChecked(w.c, w.w, xproto.ConfigWindowX, []uint32{uint32(x)})
	xproto.ConfigureWindowChecked(w.c, w.w, xproto.ConfigWindowY, []uint32{uint32(y)})
}

func (w *window) Resize(hght, wdth uint16) {
	xproto.ConfigureWindowChecked(w.c, w.w, xproto.ConfigWindowHeight, []uint32{uint32(hght)})
	xproto.ConfigureWindowChecked(w.c, w.w, xproto.ConfigWindowWidth, []uint32{uint32(wdth)})
}

func (w *window) MoveResize(x, y int16, hght, wdth uint16) {
	w.Move(x, y)
	w.Resize(hght, wdth)
}

func (w *window) Raise() {
	xproto.ConfigureWindowChecked(w.c, w.w, xproto.ConfigWindowStackMode, []uint32{xproto.StackModeAbove})
}

func (w *window) Lower() {
	xproto.ConfigureWindowChecked(w.c, w.w, xproto.ConfigWindowStackMode, []uint32{xproto.StackModeBelow})
}

func (w *window) Stack(o xproto.Window, mode uint32) {
	xproto.ConfigureWindowChecked(
		w.c,
		w.w,
		(xproto.ConfigWindowSibling | xproto.ConfigWindowStackMode),
		[]uint32{uint32(o), mode},
	)
}

func (w *window) Above(o Window) {
	w.Stack(o.XWindow(), xproto.StackModeAbove)
}

func (w *window) Below(o Window) {
	w.Stack(o.XWindow(), xproto.StackModeBelow)
}

var (
	RootEventMask uint32 = (xproto.EventMaskSubstructureNotify | xproto.EventMaskSubstructureRedirect)
	windowOff            = []uint32{RootEventMask, xproto.EventMaskSubstructureNotify} //uint32_t values_off[] = {ROOT_EVENT_MASK & ~XCB_EVENT_MASK_SUBSTRUCTURE_NOTIFY};
	windowOn             = []uint32{RootEventMask}
)

func SetVisible(v bool, c *xgb.Conn, w xproto.Window, root xproto.Window) {
	setVisibility(v, c, w, root)
}

func setVisibility(v bool, c *xgb.Conn, w xproto.Window, root xproto.Window) {
	xproto.ChangeWindowAttributesChecked(c, root, xproto.CwEventMask, windowOff)
	if v {
		xproto.MapWindow(c, w)
	} else {
		xproto.UnmapWindow(c, w)
	}
	xproto.ChangeWindowAttributesChecked(c, root, xproto.CwEventMask, windowOn)
}

func (w *window) Hide() {
	setVisibility(false, w.c, w.w, w.r)
}

func (w *window) Show() {
	setVisibility(true, w.c, w.w, w.r)
}
