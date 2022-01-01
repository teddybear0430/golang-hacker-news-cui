package ui

import (
	"fmt"
	"github.com/Yota-K/golang-hacker-news-cui/api"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"log"
	"strconv"
)

// UIの設定
func HnUi(n int) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	nodes := generateTreeNodes(n)

	t := widgets.NewTree()
	t.Title = "Hacker News CUI"
	t.TextStyle = ui.NewStyle(ui.ColorBlue)
	t.SetNodes(nodes)

	// 全画面で表示
	x, y := ui.TerminalDimensions()
	t.SetRect(0, 0, x, y)

	// キーバインディングの設定
	Keybindings(t)
}

// インターフェースを利用することで、自分好みの出力形式にカスタマイズを行なっている？
type nodeValue string

func (nv nodeValue) String() string {
	return string(nv)
}

// 取得した投稿をCUI上に表示
func generateTreeNodes(n int) []*widgets.TreeNode {
	hns := api.GetHackerNews(n)
	var nodes []*widgets.TreeNode
	for _, hn := range hns {
		node := widgets.TreeNode{
			Value: nodeValue(hn.Title),
			Nodes: []*widgets.TreeNode{
				{
					Value: nodeValue(fmt.Sprintf("Score: %s", strconv.Itoa(hn.Score))),
					Nodes: nil,
				},
				{
					Value: nodeValue(fmt.Sprintf("Type: %s", hn.Type)),
					Nodes: nil,
				},
				{
					Value: nodeValue(fmt.Sprintf("Author: %s", hn.By)),
					Nodes: nil,
				},
				{
					Value: nodeValue(fmt.Sprintf("cmd+click → %s", hn.Url)),
					Nodes: nil,
				},
				{
					Value: nodeValue("Description"),
					Nodes: []*widgets.TreeNode{
						{
							Value: nodeValue(hn.Description),
							Nodes: nil,
						},
					},
				},
			},
		}

		nodes = append(nodes, &node)
	}

	return nodes
}
