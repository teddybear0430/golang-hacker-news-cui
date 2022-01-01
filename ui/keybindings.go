package ui

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// キーバインドの設定
// switch文に登録されているキーが入力されると、それに合わせた画面描画処理と再レンダリングが行われる
func Keybindings(t *widgets.Tree) {
	uiEvents := ui.PollEvents()

	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "k":
			t.ScrollUp()
		case "j":
			t.ScrollDown()
		case "E":
			t.ExpandAll()
		case "C":
			t.CollapseAll()
		case "<Enter>":
			t.ToggleExpand()
		}

		ui.Render(t)
	}
}
