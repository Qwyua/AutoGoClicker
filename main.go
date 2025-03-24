package main

import (
	"fmt"
	"time"
	"unsafe"
	"golang.org/x/sys/windows"
)

type mouseInput struct {
	dx        int32
	dy        int32
	mouseData uint32
	dwFlags   uint32
	time      uint32
	dwExtra   uintptr
}

type input struct {
	dwType uint32
	mInput mouseInput
}

const (
	INPUT_MOUSE           = 0
	MOUSEEVENTF_LEFTDOWN  = 0x0002
	MOUSEEVENTF_LEFTUP    = 0x0004
	MOUSEEVENTF_RIGHTDOWN = 0x0008
	MOUSEEVENTF_RIGHTUP   = 0x0010
)

var (
	user32               = windows.NewLazySystemDLL("user32.dll")
	winmm                = windows.NewLazySystemDLL("winmm.dll")
	procGetAsyncKeyState = user32.NewProc("GetAsyncKeyState")
	procSendInput        = user32.NewProc("SendInput")
	procTimeBeginPeriod  = winmm.NewProc("timeBeginPeriod")
)

func main() {
	procTimeBeginPeriod.Call(uintptr(1))
	defer windows.NewLazySystemDLL("winmm.dll").NewProc("timeEndPeriod").Call(uintptr(1))

	leftCPS := 15.1  // max clicks/second
	rightCPS := 18.4 // max clicks/second

	leftInterval := time.Duration(float64(time.Second) / leftCPS)
	rightInterval := time.Duration(float64(time.Second) / rightCPS)

	var (
		nextLeft  = time.Now()
		nextRight = time.Now()
	)

	for {
		now := time.Now()

		if isKeyPressed(0x31) && now.After(nextLeft) {
			sendClick(MOUSEEVENTF_LEFTDOWN, MOUSEEVENTF_LEFTUP)
			nextLeft = now.Add(leftInterval)
		}

		if isKeyPressed(0x32) && now.After(nextRight) {
			sendClick(MOUSEEVENTF_RIGHTDOWN, MOUSEEVENTF_RIGHTUP)
			nextRight = now.Add(rightInterval)
		}

		sleepDuration := time.Until(nextLeft)
		if time.Until(nextRight) < sleepDuration {
			sleepDuration = time.Until(nextRight)
		}

		if sleepDuration > 0 {
			time.Sleep(sleepDuration)
		} else {
			time.Sleep(100 * time.Microsecond)
		}
	}
}

func isKeyPressed(vk int) bool {
	r, _, _ := procGetAsyncKeyState.Call(uintptr(vk))
	return r&0x8000 != 0
}

func sendClick(down, up uint32) {
	var i input
	i.dwType = INPUT_MOUSE
	i.mInput.dwFlags = down
	procSendInput.Call(1, uintptr(unsafe.Pointer(&i)), unsafe.Sizeof(i))
	i.mInput.dwFlags = up
	procSendInput.Call(1, uintptr(unsafe.Pointer(&i)), unsafe.Sizeof(i))
}

/* eski versiyon

package main

import (
	"time"
	"unsafe"
	"golang.org/x/sys/windows"
)

type mouseInput struct {
	dx        int32
	dy        int32
	mouseData uint32
	dwFlags   uint32
	time      uint32
	dwExtra   uintptr
}

type input struct {
	dwType uint32
	mInput mouseInput
}

const (
	INPUT_MOUSE           = 0
	MOUSEEVENTF_LEFTDOWN  = 0x0002
	MOUSEEVENTF_LEFTUP    = 0x0004
	MOUSEEVENTF_RIGHTDOWN = 0x0008
	MOUSEEVENTF_RIGHTUP   = 0x0010
)

var (
	user32               = windows.NewLazySystemDLL("user32.dll")
	procGetAsyncKeyState = user32.NewProc("GetAsyncKeyState")
	procSendInput        = user32.NewProc("SendInput")
)

func main() {
	leftCPS := 15.1  // Left Click CPS
	rightCPS := 1000 //18.4 // Right Click CPS

	leftClickInterval := time.Second / time.Duration(leftCPS)
	rightClickInterval := time.Second / time.Duration(rightCPS)

	lastLeftClickTime := time.Time{}
	lastRightClickTime := time.Time{}

	for {
		// left click icin 1 tusu
		if isKeyPressed(0x31) {
			now := time.Now()
			if now.Sub(lastLeftClickTime) >= leftClickInterval {
				sendMouseEvent(MOUSEEVENTF_LEFTDOWN)
				sendMouseEvent(MOUSEEVENTF_LEFTUP)
				lastLeftClickTime = now
			}
		}

		// right click icin 2 tusu
		if isKeyPressed(0x32) {
			now := time.Now()
			if now.Sub(lastRightClickTime) >= rightClickInterval {
				sendMouseEvent(MOUSEEVENTF_RIGHTDOWN)
				sendMouseEvent(MOUSEEVENTF_RIGHTUP)
				lastRightClickTime = now
			}
		}

		// az cpu kullanimi icin
		time.Sleep(1 * time.Millisecond)
	}
}

func isKeyPressed(vkCode int) bool {
	result, _, _ := procGetAsyncKeyState.Call(uintptr(vkCode))
	return result&0x8000 != 0
}

func sendMouseEvent(flags uint32) {
	var input input
	input.dwType = INPUT_MOUSE
	input.mInput.dwFlags = flags
	_,_,_=procSendInput.Call(uintptr(1),uintptr(unsafe.Pointer(&input)),unsafe.Sizeof(input))
}

*/
