/**
* @Author:fengxinlei
* @Description:
* @Version 1.0.0
* @Date: 2021/5/28 14:40
 */

package main
import (
	"fmt"
	termbox "github.com/nsf/termbox-go"
)
func main()  {
	fmt.Println("111111111111111")
	pause()
}

func init() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	termbox.SetCursor(0, 0)
	termbox.HideCursor()
}
func pause() {
	fmt.Println("请按任意键继续...")
Loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			break Loop
		}
	}
}
