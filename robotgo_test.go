package robotgo_test

import (
	"sync"
	"testing"

	"github.com/vcaesar/tt"
	"github.com/zuiwuchang/robotgo"
)

var m sync.Mutex

func TestMoveMouse(t *testing.T) {
	m.Lock()
	defer m.Unlock()
	robotgo.MoveMouse(20, 20)
	robotgo.MilliSleep(10)
	x, y := robotgo.GetMousePos()

	tt.Equal(t, 20, x)
	tt.Equal(t, 20, y)
}
func TestMoveMouseSmooth(t *testing.T) {
	m.Lock()
	defer m.Unlock()

	b := robotgo.MoveMouseSmooth(100, 100)
	robotgo.MilliSleep(10)
	x, y := robotgo.GetMousePos()

	tt.True(t, b)
	tt.Equal(t, 100, x)
	tt.Equal(t, 100, y)
}

func TestDragMouse(t *testing.T) {
	m.Lock()
	defer m.Unlock()

	robotgo.DragMouse(500, 500)
	robotgo.MilliSleep(10)
	x, y := robotgo.GetMousePos()

	tt.Equal(t, 500, x)
	tt.Equal(t, 500, y)
}

func TestScrollMouse(t *testing.T) {
	m.Lock()
	defer m.Unlock()

	robotgo.ScrollMouse(120, "up")
	robotgo.MilliSleep(100)

	robotgo.Scroll(210, 210)
}

func TestMoveRelative(t *testing.T) {
	m.Lock()
	defer m.Unlock()

	robotgo.Move(200, 200)
	robotgo.MilliSleep(10)

	robotgo.MoveRelative(10, -10)
	robotgo.MilliSleep(10)

	x, y := robotgo.GetMousePos()
	tt.Equal(t, 210, x)
	tt.Equal(t, 190, y)
}

func TestMoveSmoothRelative(t *testing.T) {
	m.Lock()
	defer m.Unlock()

	robotgo.Move(200, 200)
	robotgo.MilliSleep(10)

	robotgo.MoveSmoothRelative(10, -10)
	robotgo.MilliSleep(10)

	x, y := robotgo.GetMousePos()
	tt.Equal(t, 210, x)
	tt.Equal(t, 190, y)
}

func TestMouseToggle(t *testing.T) {
	e := robotgo.MouseToggle("up", "right")
	tt.Zero(t, e)
}

func TestKey(t *testing.T) {
	e := robotgo.KeyTap("v", "cmd")
	tt.Empty(t, e)

	e = robotgo.KeyToggle("v", "up")
	tt.Empty(t, e)
}

func TestClip(t *testing.T) {
	err := robotgo.WriteAll("s")
	tt.Nil(t, err)

	s, e := robotgo.ReadAll()
	tt.Equal(t, "s", s)
	tt.Nil(t, e)
}

func TestTypeStr(t *testing.T) {
	c := robotgo.CharCodeAt("s", 0)
	tt.Equal(t, 115, c)

	e := robotgo.PasteStr("s")
	tt.Empty(t, e)

	uc := robotgo.ToUC("s")
	tt.Equal(t, "[s]", uc)
}

func TestKeyCode(t *testing.T) {
	m := robotgo.MouseMap["left"]
	tt.Equal(t, 1, m)

	k := robotgo.Keycode["1"]
	tt.Equal(t, 2, k)
}
