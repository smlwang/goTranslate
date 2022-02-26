package gui

import (
	"fmt"
	"time"
	"translate/tool"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var inTE, outTE *walk.TextEdit
var selfrom, selto *walk.ComboBox
var showLB, curTime *walk.TextLabel

var showLeftFile *walk.RadioButton
var lfchecked = false

var auto *walk.RadioButton
var autoCheck = false

//暂未找到处理 ctrl + backspace 的更好办法
//var ctrl_back = ""

//字体大小
var fontSize = 12
var LangList = []string{"auto(自动)", "zh(中)", "wyw(文言)", "en(英)", "jp(日)", "kor(韩)", "th(泰)", "vie(越南)", "ru(俄)"}
var Lang = []string{"auto", "zh", "wyw", "en", "jp", "kor", "th", "vie", "ru"}

var changed = make(chan struct{}, 10)

func init() {
	go func() {
		ticker := time.Tick(500 * time.Millisecond)
		for {
			<-ticker
			curTime.SetText(fmt.Sprint(time.Now().Format("2006-01-02 15:04:05")))
		}
	}()
	go func() {
		for { //自动翻译
			<-changed
			if autoCheck {
				ok := true
				for ok {
					select {
					case <-changed:
					case <-time.After(250 * time.Millisecond):
						outTE.SetText(tool.Deal_(inTE.Text()))
						ok = false
					}
				}
			}
		}
	}()
}

//size参数凭感觉瞎调
var Mw = MainWindow{
	Title:  "简单翻译",
	Size:   Size{600, 450},
	Layout: VBox{},
	Children: []Widget{
		GroupBox{
			Layout:  HBox{},
			MaxSize: Size{100, 10},
			Children: []Widget{
				TextLabel{
					MinSize:  Size{10, 8},
					AssignTo: &showLB,
					Text:     "hello world!",
				},
				TextLabel{
					Font: Font{
						Family:    "宋体",
						PointSize: 10,
					},
					AssignTo: &curTime,
				},
			},
		},

		HSplitter{
			Children: []Widget{
				TextEdit{
					MinSize:  Size{180, 50},
					VScroll:  true,
					AssignTo: &inTE,
					Font: Font{
						PointSize: fontSize,
					},
					OnKeyDown: func(key walk.Key) {
						if key == walk.KeyBack { //防止对空白文本框狂操退格键使程序退出
							if inTE.Text() == "" {
								return
							}
						}
					},
					OnTextChanged: func() {
						txt := inTE.Text()
						if autoCheck && txt != "" {
							changed <- struct{}{}
						}
					},
				},
				GroupBox{
					MaxSize:   Size{20, 0},
					Layout:    VBox{},
					Alignment: AlignHCenterVCenter,
					RowSpan:   0,
					Children: []Widget{
						ComboBox{
							Alignment:       Alignment2D(AlignCenter),
							AssignTo:        &selfrom,
							ToolTipText:     "源语言",
							Model:           LangList,
							CurrentIndex:    0,
							DoubleBuffering: true,
							OnCurrentIndexChanged: func() {
								tool.SetFrom(Lang[selfrom.CurrentIndex()])
							},
						},
						TextLabel{
							MaxSize: Size{15, 5},
							Text:    "\t\t||\n\t\t||\n\t\t||\n\t           \\--/\n\t\t\\/",
						},
						ComboBox{
							Alignment:    Alignment2D(AlignCenter),
							AssignTo:     &selto,
							ToolTipText:  "目标语言",
							Model:        LangList[1:],
							CurrentIndex: 0,
							OnCurrentIndexChanged: func() {
								idx := selto.CurrentIndex()
								tool.SetTo(Lang[idx+1])
							},
						},
						PushButton{
							ToolTipText: "导入文本翻译",
							Text:        "打开(导入)",
							OnClicked:   selectFile,
						},
						RadioButton{
							AssignTo: &showLeftFile,
							Text:     "将(导入)文本显示在左侧",
							OnClicked: func() {
								showLeftFile.SetChecked(!lfchecked)
								lfchecked = !lfchecked
							},
						},
						RadioButton{
							AssignTo: &auto,
							Text:     "自动翻译",
							OnClicked: func() {
								auto.SetChecked(!autoCheck)
								autoCheck = !autoCheck
							},
						},
					},
				},
				TextEdit{
					VScroll:  true,
					MinSize:  Size{190, 50},
					AssignTo: &outTE,
					ReadOnly: true,
					Font: Font{
						PointSize: fontSize,
					},
				},
			},
		},
		GroupBox{
			Layout:  HBox{},
			MaxSize: Size{0, 50},
			Children: []Widget{
				GroupBox{
					MaxSize:    Size{100, 10},
					Layout:     HBox{},
					ColumnSpan: 5,
					Column:     1,
					Children: []Widget{
						ToolButton{
							Text: "清空",
							OnClicked: func() {
								inTE.SetText("")
							},
							MinSize: Size{30, 5},
						},
						ToolButton{
							Text: "粘贴",
							OnClicked: func() {
								txt, err := walk.Clipboard().Text()
								if err != nil {
									showLB.SetText("faild when paste")
								} else {
									inTE.SetText(txt)
									showLB.SetText("paste ok")
								}
							},
							MinSize: Size{30, 5},
						},
					},
				},
				PushButton{
					MinSize: Size{100, 50},
					Text:    "翻译",
					OnClicked: func() {
						outTE.SetText(tool.Deal_(inTE.Text()))
					},
				},
				GroupBox{
					MaxSize:    Size{100, 10},
					ColumnSpan: 5,
					Column:     1,
					Layout:     HBox{},
					Children: []Widget{
						ToolButton{
							MinSize: Size{30, 5},
							Text:    "复制",
							OnClicked: func() {
								if err := walk.Clipboard().SetText(outTE.Text()); err != nil {
									showLB.SetText("faild when copy")
								} else {
									showLB.SetText("copy ok")
								}

							},
						},
					},
				},
			},
		},
	},
}

func selectFile() {
	dlg := walk.FileDialog{
		Title:  "选择文件",
		Filter: " 文本文件 (*.txt)| *.txt| 所有文件(*.*)| *.*",
	}
	if ok, err := dlg.ShowOpen(&walk.MainWindow{}); err != nil {
		showLB.SetText(err.Error())
		return
	} else if !ok {
		showLB.SetText("faild to open")
		return
	}
	if m, txt, ok := tool.Deal(dlg.FilePath); ok {
		if txt[0] == '\n' {
			txt = txt[1:]
		}
		outTE.SetText(txt)
		if lfchecked {
			inTE.AppendText(m)
		}
		showLB.SetText("success!")
	} else {
		showLB.SetText(txt)
	}
}
