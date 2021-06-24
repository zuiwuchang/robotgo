package robotgo

/*
//#if defined(IS_MACOSX)
	#cgo darwin CFLAGS: -x objective-c -Wno-deprecated-declarations
	#cgo darwin LDFLAGS: -framework Cocoa -framework OpenGL -framework IOKit
	#cgo darwin LDFLAGS: -framework Carbon -framework CoreFoundation
//#elif defined(USE_X11)
	// Drop -std=c11
	#cgo linux CFLAGS: -I/usr/src
	#cgo linux LDFLAGS: -L/usr/src -lpng -lz -lX11 -lXtst -lm
	// #cgo linux LDFLAGS: -lX11-xcb -lxcb -lxcb-xkb -lxkbcommon -lxkbcommon-x11
//#endif
	#cgo windows LDFLAGS: -lgdi32 -luser32
#include "screen/goScreen.h"
#include "mouse/goMouse.h"
#include "key/goKey.h"
*/
import "C"
import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/zuiwuchang/robotgo/clipboard"

	"github.com/vcaesar/tt"
)

// MilliSleep sleep tm milli second
func MilliSleep(tm int) {
	time.Sleep(time.Duration(tm) * time.Millisecond)
}

// Sleep time.Sleep tm second
func Sleep(tm int) {
	time.Sleep(time.Duration(tm) * time.Second)
}

// MicroSleep time C.microsleep(tm)
func MicroSleep(tm float64) {
	C.microsleep(C.double(tm))
}

/*
.___  ___.   ______    __    __       _______. _______
|   \/   |  /  __  \  |  |  |  |     /       ||   ____|
|  \  /  | |  |  |  | |  |  |  |    |   (----`|  |__
|  |\/|  | |  |  |  | |  |  |  |     \   \    |   __|
|  |  |  | |  `--'  | |  `--'  | .----)   |   |  |____
|__|  |__|  \______/   \______/  |_______/    |_______|

*/

// CheckMouse check the mouse button
func CheckMouse(btn string) C.MMMouseButton {
	// button = args[0].(C.MMMouseButton)
	if btn == "left" {
		return C.LEFT_BUTTON
	}

	if btn == "center" {
		return C.CENTER_BUTTON
	}

	if btn == "right" {
		return C.RIGHT_BUTTON
	}

	return C.LEFT_BUTTON
}

// MoveMouse move the mouse
func MoveMouse(x, y int) {
	// C.size_t  int
	Move(x, y)
}

// Move move the mouse
func Move(x, y int) {
	cx := C.int32_t(x)
	cy := C.int32_t(y)
	C.move_mouse(cx, cy)
}

// DragMouse drag the mouse
func DragMouse(x, y int, args ...string) {
	Drag(x, y, args...)
}

// Drag drag the mouse
func Drag(x, y int, args ...string) {
	var button C.MMMouseButton = C.LEFT_BUTTON

	cx := C.int32_t(x)
	cy := C.int32_t(y)

	if len(args) > 0 {
		button = CheckMouse(args[0])
	}

	C.drag_mouse(cx, cy, button)
}

// DragSmooth drag the mouse smooth
func DragSmooth(x, y int, args ...interface{}) {
	MouseToggle("down")
	MoveSmooth(x, y, args...)
	MouseToggle("up")
}

// MoveMouseSmooth move the mouse smooth,
// moves mouse to x, y human like, with the mouse button up.
func MoveMouseSmooth(x, y int, args ...interface{}) bool {
	return MoveSmooth(x, y, args...)
}

// MoveSmooth move the mouse smooth,
// moves mouse to x, y human like, with the mouse button up.
//
// robotgo.MoveSmooth(x, y int, low, high float64, mouseDelay int)
func MoveSmooth(x, y int, args ...interface{}) bool {
	cx := C.int32_t(x)
	cy := C.int32_t(y)

	var (
		mouseDelay = 10
		low        C.double
		high       C.double
	)

	if len(args) > 2 {
		mouseDelay = args[2].(int)
	}

	if len(args) > 1 {
		low = C.double(args[0].(float64))
		high = C.double(args[1].(float64))
	} else {
		low = 1.0
		high = 3.0
	}

	cbool := C.move_mouse_smooth(cx, cy, low, high, C.int(mouseDelay))

	return bool(cbool)
}

// MoveArgs move mose relative args
func MoveArgs(x, y int) (int, int) {
	mx, my := GetMousePos()
	mx = mx + x
	my = my + y

	return mx, my
}

// MoveRelative move mose relative
func MoveRelative(x, y int) {
	Move(MoveArgs(x, y))
}

// MoveSmoothRelative move mose smooth relative
func MoveSmoothRelative(x, y int, args ...interface{}) {
	mx, my := MoveArgs(x, y)
	MoveSmooth(mx, my, args...)
}

// GetMousePos get mouse's portion
func GetMousePos() (int, int) {
	pos := C.get_mouse_pos()

	x := int(pos.x)
	y := int(pos.y)

	return x, y
}

// MouseClick click the mouse
//
// robotgo.MouseClick(button string, double bool)
func MouseClick(args ...interface{}) {
	Click(args...)
}

// Click click the mouse
//
// robotgo.Click(button string, double bool)
func Click(args ...interface{}) {
	var (
		button C.MMMouseButton = C.LEFT_BUTTON
		double C.bool
	)

	if len(args) > 0 {
		button = CheckMouse(args[0].(string))
	}

	if len(args) > 1 {
		double = C.bool(args[1].(bool))
	}

	C.mouse_click(button, double)
}

// MoveClick move and click the mouse
//
// robotgo.MoveClick(x, y int, button string, double bool)
func MoveClick(x, y int, args ...interface{}) {
	MoveMouse(x, y)
	MouseClick(args...)
}

// MovesClick move smooth and click the mouse
func MovesClick(x, y int, args ...interface{}) {
	MoveSmooth(x, y)
	MouseClick(args...)
}

// MouseToggle toggle the mouse
func MouseToggle(togKey string, args ...interface{}) int {
	var button C.MMMouseButton = C.LEFT_BUTTON

	if len(args) > 0 {
		button = CheckMouse(args[0].(string))
	}

	down := C.CString(togKey)
	i := C.mouse_toggle(down, button)

	C.free(unsafe.Pointer(down))
	return int(i)
}

// ScrollMouse scroll the mouse
func ScrollMouse(x int, direction string) {
	cx := C.size_t(x)
	cy := C.CString(direction)
	C.scroll_mouse(cx, cy)

	C.free(unsafe.Pointer(cy))
}

// Scroll scroll the mouse with x, y
//
// robotgo.Scroll(x, y, msDelay int)
func Scroll(x, y int, args ...int) {
	var msDelay = 10
	if len(args) > 0 {
		msDelay = args[0]
	}

	cx := C.int(x)
	cy := C.int(y)
	cz := C.int(msDelay)

	C.scroll(cx, cy, cz)
}

// SetMouseDelay set mouse delay
func SetMouseDelay(delay int) {
	cdelay := C.size_t(delay)
	C.set_mouse_delay(cdelay)
}

/*
 __  ___  ___________    ____ .______     ______        ___      .______       _______
|  |/  / |   ____\   \  /   / |   _  \   /  __  \      /   \     |   _  \     |       \
|  '  /  |  |__   \   \/   /  |  |_)  | |  |  |  |    /  ^  \    |  |_)  |    |  .--.  |
|    <   |   __|   \_    _/   |   _  <  |  |  |  |   /  /_\  \   |      /     |  |  |  |
|  .  \  |  |____    |  |     |  |_)  | |  `--'  |  /  _____  \  |  |\  \----.|  '--'  |
|__|\__\ |_______|   |__|     |______/   \______/  /__/     \__\ | _| `._____||_______/

*/

// KeyTap tap the keyboard code;
//
// See keys:
//	https://github.com/zuiwuchang/robotgo/blob/master/docs/keys.md
//
func KeyTap(tapKey string, args ...interface{}) string {
	var (
		akey     string
		keyT     = "null"
		keyArr   []string
		num      int
		keyDelay = 10
	)

	// var ckeyArr []*C.char
	ckeyArr := make([](*C.char), 0)
	// zkey := C.CString(args[0])
	zkey := C.CString(tapKey)
	defer C.free(unsafe.Pointer(zkey))

	if len(args) > 2 && (reflect.TypeOf(args[2]) != reflect.TypeOf(num)) {
		num = len(args)
		for i := 0; i < num; i++ {
			s := args[i].(string)
			ckeyArr = append(ckeyArr, (*C.char)(unsafe.Pointer(C.CString(s))))
		}

		str := C.key_Taps(zkey,
			(**C.char)(unsafe.Pointer(&ckeyArr[0])), C.int(num), 0)
		return C.GoString(str)
	}

	if len(args) > 0 {
		if reflect.TypeOf(args[0]) == reflect.TypeOf(keyArr) {

			keyArr = args[0].([]string)
			num = len(keyArr)
			for i := 0; i < num; i++ {
				ckeyArr = append(ckeyArr, (*C.char)(unsafe.Pointer(C.CString(keyArr[i]))))
			}

			if len(args) > 1 {
				keyDelay = args[1].(int)
			}
		} else {
			akey = args[0].(string)

			if len(args) > 1 {
				if reflect.TypeOf(args[1]) == reflect.TypeOf(akey) {
					keyT = args[1].(string)
					if len(args) > 2 {
						keyDelay = args[2].(int)
					}
				} else {
					keyDelay = args[1].(int)
				}
			}
		}

	} else {
		akey = "null"
		keyArr = []string{"null"}
	}

	if akey == "" && len(keyArr) != 0 {
		str := C.key_Taps(zkey, (**C.char)(unsafe.Pointer(&ckeyArr[0])),
			C.int(num), C.int(keyDelay))

		return C.GoString(str)
	}

	amod := C.CString(akey)
	amodt := C.CString(keyT)
	str := C.key_tap(zkey, amod, amodt, C.int(keyDelay))

	C.free(unsafe.Pointer(amod))
	C.free(unsafe.Pointer(amodt))

	return C.GoString(str)
}

// KeyToggle toggle the keyboard
//
// See keys:
//	https://github.com/zuiwuchang/robotgo/blob/master/docs/keys.md
//
func KeyToggle(key string, args ...string) string {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	ckeyArr := make([](*C.char), 0)
	if len(args) > 3 {
		num := len(args)
		for i := 0; i < num; i++ {
			ckeyArr = append(ckeyArr, (*C.char)(unsafe.Pointer(C.CString(args[i]))))
		}

		str := C.key_Toggles(ckey, (**C.char)(unsafe.Pointer(&ckeyArr[0])), C.int(num))
		return C.GoString(str)
	}

	var (
		down, mKey, mKeyT = "null", "null", "null"
		// keyDelay = 10
	)

	if len(args) > 0 {
		down = args[0]

		if len(args) > 1 {
			mKey = args[1]
			if len(args) > 2 {
				mKeyT = args[2]
			}
		}
	}

	cdown := C.CString(down)
	cmKey := C.CString(mKey)
	cmKeyT := C.CString(mKeyT)

	str := C.key_toggle(ckey, cdown, cmKey, cmKeyT)
	// str := C.key_Toggle(ckey, cdown, cmKey, cmKeyT, C.int(keyDelay))

	C.free(unsafe.Pointer(cdown))
	C.free(unsafe.Pointer(cmKey))
	C.free(unsafe.Pointer(cmKeyT))

	return C.GoString(str)
}

// ReadAll read string from clipboard
func ReadAll() (string, error) {
	return clipboard.ReadAll()
}

// WriteAll write string to clipboard
func WriteAll(text string) error {
	return clipboard.WriteAll(text)
}

// CharCodeAt char code at utf-8
func CharCodeAt(s string, n int) rune {
	i := 0
	for _, r := range s {
		if i == n {
			return r
		}
		i++
	}

	return 0
}

// UnicodeType tap uint32 unicode
func UnicodeType(str uint32) {
	cstr := C.uint(str)
	C.unicodeType(cstr)
}
func ToUC(text string) []string {
	return toUC(text)
}
func toUC(text string) []string {
	var uc []string

	for _, r := range text {
		textQ := strconv.QuoteToASCII(string(r))
		textUnQ := textQ[1 : len(textQ)-1]

		st := strings.Replace(textUnQ, "\\u", "U", -1)
		uc = append(uc, st)
	}

	return uc
}

func inputUTF(str string) {
	cstr := C.CString(str)
	C.input_utf(cstr)

	C.free(unsafe.Pointer(cstr))
}

// TypeStr send a string, support UTF-8
//
// robotgo.TypeStr(string: The string to send, float64: microsleep time, x11)
func TypeStr(str string, args ...float64) {
	var tm, tm1 = 0.0, 7.0

	if len(args) > 0 {
		tm = args[0]
	}
	if len(args) > 1 {
		tm1 = args[1]
	}

	if runtime.GOOS == "linux" {
		strUc := toUC(str)
		for i := 0; i < len(strUc); i++ {
			ru := []rune(strUc[i])
			if len(ru) <= 1 {
				ustr := uint32(CharCodeAt(strUc[i], 0))
				UnicodeType(ustr)
			} else {
				inputUTF(strUc[i])
				MicroSleep(tm1)
			}

			MicroSleep(tm)
		}

		return
	}

	for i := 0; i < len([]rune(str)); i++ {
		ustr := uint32(CharCodeAt(str, i))
		UnicodeType(ustr)

		// if len(args) > 0 {
		MicroSleep(tm)
		// }
	}
}

// PasteStr paste a string, support UTF-8
func PasteStr(str string) string {
	err := clipboard.WriteAll(str)
	if err != nil {
		return fmt.Sprint(err)
	}

	if runtime.GOOS == "darwin" {
		return KeyTap("v", "command")
	}

	return KeyTap("v", "control")
}

// TypeString send a string, support unicode
// TypeStr(string: The string to send), Wno-deprecated
func TypeString(str string, delay ...int) {
	tt.Drop("TypeString", "TypeStr")
	var cdelay C.size_t
	cstr := C.CString(str)
	if len(delay) > 0 {
		cdelay = C.size_t(delay[0])
	}

	C.type_string_delayed(cstr, cdelay)

	C.free(unsafe.Pointer(cstr))
}

// TypeStrDelay type string delayed
func TypeStrDelay(str string, delay int) {
	TypeStr(str)
	Sleep(delay)
}

// TypeStringDelayed type string delayed, Wno-deprecated
func TypeStringDelayed(str string, delay int) {
	tt.Drop("TypeStringDelayed", "TypeStrDelay")
	TypeStrDelay(str, delay)
}

// SetKeyDelay set keyboard delay
func SetKeyDelay(delay int) {
	C.set_keyboard_delay(C.size_t(delay))
}

// SetKeyboardDelay set keyboard delay, Wno-deprecated,
// this function will be removed in version v1.0.0
func SetKeyboardDelay(delay int) {
	tt.Drop("SetKeyboardDelay", "SetKeyDelay")
	SetKeyDelay(delay)
}

// SetDelay set the key and mouse delay
func SetDelay(d ...int) {
	v := 10
	if len(d) > 0 {
		v = d[0]
	}

	SetMouseDelay(v)
	SetKeyDelay(v)
}
