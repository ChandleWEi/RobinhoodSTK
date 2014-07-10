package main

import (
	"../ARHelper"
	"fmt"
	// "github.com/salviati/go-qt5/qt5"
	"github.com/ChandleWEi/go-qt5/qt5"
	"image/color"
	"io/ioutil"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"runtime"
	"strconv"
	"time"
)

const (
	originType = 0
	macdType   = 1
)

type rgba uint32

var gWin *MainWindow

func (c rgba) RGBA() (r, g, b, a uint32) {
	return uint32((c >> 16) & 0xff), uint32((c >> 8) & 0xff), uint32(c & 0xff), uint32(c >> 24)
}

func make_rgba(c color.Color) rgba {
	if c == nil {
		return 0
	}
	r, g, b, a := c.RGBA()
	return rgba(((a & 0xff) << 24) | ((r & 0xff) << 16) | ((g & 0xff) << 8) | (b & 0xff))
}

func main() {

	qt5.Main(ui_main)
}

func ui_main() {
	exit := make(chan bool)
	go func() {
		fmt.Println("vfc/ui")
		qt5.OnInsertObject(func(v interface{}) {
			fmt.Println("add item", v)
		})
		qt5.OnRemoveObject(func(v interface{}) {
			fmt.Println("remove item", v)
		})
		w := new(MainWindow).Init()
		gWin = w
		defer w.Close()

		w.SetSizev(800, 600)
		w.OnCloseEvent(func(e *qt5.CloseEvent) {
			fmt.Println("close", e)
		})
		w.Show()
		<-exit
	}()
	qt5.Run()
	exit <- true
}

type MainWindow struct {
	qt5.Widget
	tab  *qt5.TabWidget
	sbar *qt5.StatusBar
}

func (p *MainWindow) createStockTab() *qt5.Widget {
	w := qt5.NewWidget()
	vbox := qt5.NewVBoxLayout()
	hbox := qt5.NewHBoxLayout()
	my := new(MyWidget).StockInit()
	lbl := qt5.NewLabel()
	lbl.SetText("stock macd example - draw lines")

	btn := qt5.NewButton()
	btn.SetText("Load")
	btn.OnClicked(func() {
		my.Load()
	})

	btn1 := qt5.NewButton()
	btn1.SetText("Clear")
	btn1.OnClicked(func() {
		my.Clear()
	})

	btn2 := qt5.NewButton()
	btn2.SetText("1.下载上海列表CSV")
	btn2.OnClicked(func() {
		ARHelper.DownloadShStockCsv()
	})

	btn3 := qt5.NewButton()
	btn3.SetText("2.下载上海列表数据CSV")
	btn3.OnClicked(func() {
		ARHelper.DownloadAllCSV()
	})

	btn4 := qt5.NewButton()
	btn4.SetText("3.导入数据到mgo")
	btn4.OnClicked(func() {
		ARHelper.BuildAll()
	})

	hbox.AddWidget(lbl)
	hbox.AddWidgetWith(btn, 0, qt5.AlignRight)
	hbox.AddWidgetWith(btn1, 0, qt5.AlignRight)
	hbox.AddWidgetWith(btn2, 0, qt5.AlignRight)
	hbox.AddWidgetWith(btn3, 0, qt5.AlignRight)
	hbox.AddWidgetWith(btn4, 0, qt5.AlignRight)
	vbox.AddLayout(hbox)
	vbox.AddWidgetWith(my, 1, 0)
	w.SetLayout(vbox)
	return w
}

func (p *MainWindow) Init() *MainWindow {
	if p.Widget.Init() == nil {
		return nil
	}
	p.SetWindowTitle("MainWindow")

	p.tab = qt5.NewTabWidget()

	p.tab.AddTab(p.createStockTab(), "Stock", nil)
	// p.tab.AddTab(p.createStdTab(), "Standard", nil)
	// p.tab.AddTab(p.createMyTab(), "Custom", nil)
	// p.tab.AddTab(p.createToolBox(), "ToolBox", nil)

	p.sbar = qt5.NewStatusBar()

	menubar := qt5.NewMenuBar()
	menu := qt5.NewMenuWithTitle("&File")
	//menu.SetTitle("&File")
	menubar.AddMenu(menu)

	act := qt5.NewAction()
	act.SetText("&Quit")
	act.OnTriggered(func(bool) {
		p.Close()
	})
	ic := qt5.NewIconWithFile("images/close.png")
	//defer ic.Close()
	act.SetIcon(ic)
	menu.AddAction(act)

	toolBar := qt5.NewToolBar()
	toolBar.AddAction(act)
	toolBar.AddSeparator()
	cmb := qt5.NewComboBox()
	cmb.AddItem("600705")
	cmb.AddItem("600000")
	cmb.SetToolTip("ComboBox")
	cmbAct := toolBar.AddWidget(cmb)

	fm := qt5.NewComboBox()
	fm.AddItem("macd")
	fm.AddItem("kdj1")
	fm.AddItem("kdj2")
	fm.AddItem("kdj3")
	fm.SetCurrentIndex(0)
	fm.SetToolTip("formula")
	toolBar.AddWidget(fm)

	fmt.Println(cmbAct)

	vbox := qt5.NewVBoxLayout()
	vbox.SetMargin(0)
	vbox.SetSpacing(0)
	vbox.SetMenuBar(menubar)
	vbox.AddWidget(toolBar)
	vbox.AddWidget(p.tab)
	vbox.AddWidget(p.sbar)

	p.SetLayout(vbox)

	p.tab.OnCurrentChanged(func(index int) {
		p.sbar.ShowMessage("current: "+p.tab.TabText(index), 0)
	})

	systray := qt5.NewSystemTray()
	systray.SetContextMenu(menu)
	systray.SetIcon(ic)
	systray.SetVisible(true)
	systray.ShowMessage("hello", "this is a test", qt5.Information, 1000)
	ic2 := systray.Icon()
	fmt.Println(ic2)

	p.SetWindowIcon(ic2)

	return p
}

func (p *MainWindow) createStdTab() *qt5.Widget {
	w := qt5.NewWidget()
	vbox := qt5.NewVBoxLayout()
	w.SetLayout(vbox)

	ed := qt5.NewLineEdit()
	ed.SetInputMask("0000-00-00")
	ed.SetText("2012-01-12")

	lbl := qt5.NewLabel()
	lbl.SetText("Label")
	btn := qt5.NewButton()
	btn.SetText("Button")
	chk := qt5.NewCheckBox()
	chk.SetText("CheckBox")
	radio := qt5.NewRadio()
	radio.SetText("Radio")
	cmb := qt5.NewComboBox()
	cmb.AddItem("001")
	cmb.AddItem("002")
	cmb.AddItem("003")
	cmb.SetCurrentIndex(2)
	fmt.Println(cmb.CurrentIndex())
	cmb.OnCurrentIndexChanged(func(v int) {
		fmt.Println(cmb.ItemText(v))
	})

	slider := qt5.NewSlider()
	slider.SetTickInterval(50)
	slider.SetTickPosition(qt5.TicksBothSides)
	slider.SetSingleStep(1)

	scl := qt5.NewScrollBar()
	fmt.Println(slider.Range())

	dial := qt5.NewDial()

	dial.SetNotchesVisible(true)
	dial.SetNotchTarget(10)
	fmt.Println(dial.NotchSize())

	vbox.AddWidget(ed)
	vbox.AddWidget(lbl)
	vbox.AddWidget(btn)
	vbox.AddWidget(chk)
	vbox.AddWidget(radio)
	vbox.AddWidget(cmb)
	vbox.AddWidget(slider)
	vbox.AddWidget(scl)
	vbox.AddWidget(dial)
	vbox.AddStretch(0)
	return w
}

func (p *MainWindow) createToolBox() qt5.IWidget {
	tb := qt5.NewToolBox()
	tb.AddItem(qt5.NewButtonWithText("button"), "btn", nil)
	tb.AddItem(qt5.NewLabelWithText("Label\nInfo"), "Label", nil)
	pixmap := qt5.NewPixmapWithFile("images/liteide128.png")
	//defer pixmap.Close()
	lbl := qt5.NewLabel()
	lbl.SetPixmap(pixmap)
	tb.AddItem(lbl, "Lalel Pixmap", nil)
	buf, err := ioutil.ReadFile("images/liteide128.png")
	if err == nil {
		pixmap2 := qt5.NewPixmapWithData(buf)
		tb.AddItem(qt5.NewLabelWithPixmap(pixmap2), "Lalel Pixmap2", nil)
	}
	return tb
}

func (p *MainWindow) createMyTab() *qt5.Widget {
	w := qt5.NewWidget()
	vbox := qt5.NewVBoxLayout()
	hbox := qt5.NewHBoxLayout()
	my := new(MyWidget).Init()
	lbl := qt5.NewLabel()
	lbl.SetText("this is custome widget - draw lines")
	btn := qt5.NewButton()
	btn.SetText("Clear")
	btn.OnClicked(func() {
		my.Clear()
	})
	hbox.AddWidget(lbl)
	hbox.AddWidgetWith(btn, 0, qt5.AlignRight)
	vbox.AddLayout(hbox)
	vbox.AddWidgetWith(my, 1, 0)
	w.SetLayout(vbox)
	return w
}

type MyWidget struct {
	qt5.Widget
	lines  [][]qt5.Point
	line   []qt5.Point
	font   *qt5.Font
	result *[]map[string]interface{}
}

func (p *MyWidget) Name() string {
	return "MyWidget"
}

func (p *MyWidget) String() string {
	return qt5.DumpObject(p)
}

func (p *MyWidget) StockInit() *MyWidget {
	if p.Widget.Init() == nil {
		return nil
	}
	p.font = qt5.NewFontWith("Timer", 16, 87)
	p.font.SetItalic(true)
	p.Widget.OnPaintEvent(func(e *qt5.PaintEvent) {
		p.stockPaintEvent(e)
	})
	// p.Widget.OnMousePressEvent(func(e *qt5.MouseEvent) {
	// 	p.mousePressEvent(e)
	// })
	// p.Widget.OnMouseMoveEvent(func(e *qt5.MouseEvent) {
	// 	p.mouseMoveEvent(e)
	// })
	// p.Widget.OnMouseReleaseEvent(func(e *qt5.MouseEvent) {
	// 	p.mouseReleaseEvent(e)
	// })
	// qt5.InsertObject(p)
	return p
}

func (p *MyWidget) stockPaintEvent(e *qt5.PaintEvent) {
	if p.result != nil {

		// draw background coordinate
		paint := qt5.NewPainter()
		defer paint.Close()
		paint.Begin(p)
		font := qt5.NewFontWith("Timer", 16, 87)
		font.SetItalic(true)
		paint.SetFont(font)
		x, y := p.Widget.Sizev()
		for i := x; i > 0; i-- {
			t := strconv.Itoa(i)
			text := t + ",0"
			if i%60 == 0 {
				paint.DrawText(qt5.Pt(i, y), text)
			}
		}

		for i := y; i > 0; i-- {
			t := strconv.Itoa(i)
			text := "0," + t
			if i%20 == 0 {
				paint.DrawText(qt5.Pt(0, y-i), text)
			}
		}

		pen9 := qt5.NewPen()
		pen9.SetColor(color.RGBA{255, 128, 0, 0})
		pen9.SetWidth(2)
		painterDrawLine(p.result, pen9, &p.Widget, "ema9")

		pen12 := qt5.NewPen()
		pen12.SetColor(color.RGBA{89, 77, 76, 0})
		pen12.SetWidth(2)
		pen12.SetStyle(qt5.SolidLine)
		painterDrawLine(p.result, pen12, &p.Widget, "ema12")

		pen26 := qt5.NewPen()
		pen26.SetColor(color.RGBA{255, 61, 46, 0})
		pen26.SetWidth(2)
		pen26.SetStyle(qt5.SolidLine)
		painterDrawLine(p.result, pen26, &p.Widget, "ema26")
	}

	// paint.Begin(p)
	// paint.SetFont(p.font)
	// paint.DrawLines(p.line)
	// paint.SetFont(p.font)
	// paint.DrawText(qt5.Pt(100, 100), "draw test")
	// for _, v := range p.lines {
	// 	//paint.DrawLines(v)
	// 	paint.DrawPolyline(v)
	// }
	// paint.End()
	// runtime.GC()
}

func (p *MyWidget) Load() {
	p.lines = [][]qt5.Point{}
	result := getMacdData()
	p.result = &result
	p.Update()
}

func (p *MyWidget) Clear() {
	p.lines = [][]qt5.Point{}
	p.result = nil
	p.Update()
}

func (p *MyWidget) Init() *MyWidget {
	if p.Widget.Init() == nil {
		return nil
	}
	p.font = qt5.NewFontWith("Timer", 16, 87)
	p.font.SetItalic(true)
	p.Widget.OnPaintEvent(func(e *qt5.PaintEvent) {
		p.paintEvent(e)
	})
	p.Widget.OnMousePressEvent(func(e *qt5.MouseEvent) {
		p.mousePressEvent(e)
	})
	p.Widget.OnMouseMoveEvent(func(e *qt5.MouseEvent) {
		p.mouseMoveEvent(e)
	})
	p.Widget.OnMouseReleaseEvent(func(e *qt5.MouseEvent) {
		p.mouseReleaseEvent(e)
	})
	qt5.InsertObject(p)
	return p
}

func (p *MyWidget) paintEvent(e *qt5.PaintEvent) {
	paint := qt5.NewPainter()
	defer paint.Close()

	paint.Begin(p)
	paint.SetFont(p.font)
	paint.DrawLines(p.line)
	paint.SetFont(p.font)
	paint.DrawText(qt5.Pt(100, 100), "draw test")
	for _, v := range p.lines {
		//paint.DrawLines(v)
		paint.DrawPolyline(v)
	}
	paint.End()
	runtime.GC()
}

func (p *MyWidget) mousePressEvent(e *qt5.MouseEvent) {
	p.line = append(p.line, e.Pos())
	p.Update()
}

func (p *MyWidget) mouseMoveEvent(e *qt5.MouseEvent) {
	p.line = append(p.line, e.Pos())
	p.Update()
}

func (p *MyWidget) mouseReleaseEvent(e *qt5.MouseEvent) {
	p.line = append(p.line, e.Pos())
	p.lines = append(p.lines, p.line)
	p.line = []qt5.Point{}
	p.Update()
}

// func oldui_main() {
// 	exit := make(chan bool)
// 	// 	go func() {
// 	// 		go ui_main()
// 	// 		//		qt5.SetStyleSheet("QLineEdit { background-color: yellow }")
// 	// 		qt5.App().SetStyleSheet(`
// 	//  QFrame, QLabel, QToolTip {
// 	//      border: 2px solid green;
// 	//      border-radius: 4px;
// 	//      padding: 2px;
// 	//      background-color: yellow;
// 	//  }
// 	// `)
// 	// 		// background-image: url(images/welcome.png);
// 	// 		fmt.Println("font is ", qt5.App().StyleSheet())
// 	// 		qt5.Run()
// 	// 		exit <- true
// 	// 	}

// 	fmt.Println("ui_main is runing ")
// 	w := qt5.NewWidget()

// 	// vbox := qt5.NewVBoxLayout()
// 	// w.SetLayout(vbox)
// 	// lbl := qt5.NewLabel()
// 	// lbl.SetText("<h2><i>Hello</i> <font color=blue><a href=\"ui\">UI</a></font></h2>")
// 	//	lbl.OnLinkActivated(fnTEST)
// 	// vbox.AddWidget(lbl)
// 	// vbox.AddStretch(0)
// 	result := getMacdData()

// 	w.OnPaintEvent(func(e *qt5.PaintEvent) {

// 		pen9 := qt5.NewPen()
// 		pen9.SetColor(color.RGBA{255, 128, 0, 0})
// 		pen9.SetWidth(2)
// 		painterDrawLine(&result, pen9, w, "ema9")

// 		pen12 := qt5.NewPen()
// 		pen12.SetColor(color.RGBA{89, 77, 76, 0})
// 		pen12.SetWidth(2)
// 		pen12.SetStyle(qt5.SolidLine)
// 		painterDrawLine(&result, pen12, w, "ema12")

// 		pen26 := qt5.NewPen()
// 		pen26.SetColor(color.RGBA{255, 61, 46, 0})
// 		pen26.SetWidth(2)
// 		pen26.SetStyle(qt5.SolidLine)
// 		painterDrawLine(&result, pen26, w, "ema26")

// 	})
// 	// qt5.Version()
// 	w.SetWindowTitle("stock 0.0.1")
// 	w.SetSizev(500, 500)
// 	defer w.Close()
// 	w.Show()
// 	<-exit
// }

func painterDrawLine(result *[]map[string]interface{}, pen *qt5.Pen, w *qt5.Widget, key string) {

	var startPt qt5.Point

	var endPt qt5.Point
	paint := qt5.NewPainter()
	defer paint.Close()

	paint.Begin(w)
	for i, v := range *result {
		if i == 0 {
			_, y := gWin.Widget.Sizev()
			startPt = qt5.Pt(i*5, (y - int(v[key].(float64))*10))

		} else {

			startPt = endPt
		}
		x, y := gWin.Widget.Sizev()
		fmt.Println("Widget x is ", x, " y is ", y)
		endPt = qt5.Pt(i*5, (y - int(v[key].(float64))*10))
		paint.SetPen(pen)
		paint.DrawLine(startPt, endPt)

	}

	paint.End()
	runtime.GC()
}

// mongo 抓取所有macd 数据
func getMacdData() (result []map[string]interface{}) {

	var (
		session *mgo.Session
		db      *mgo.Database
		err     error
	)

	// init mongo
	//	session, err := mgo.Dial("localhost")
	if session, err = mgo.Dial("localhost"); err != nil {
		panic(err)
	}
	db = session.DB("robinhood")

	collString := "ss600705"
	collection := db.C(collString)
	// collCount, _ := collection.Count()
	// fmt.Println("collection name is ", collection.Name, " count is ", collCount)
	start := time.Now().Add(-24 * 300 * time.Hour)
	// query := map[string]map[string]time.Time{"date": {"$gte": start}}
	// query := map[string]interface{}{"type": originType, "date": {"$gte": start}}
	query := bson.M{"type": macdType, "date": bson.M{"$gte": start}}
	//			}
	queryRS := collection.Find(query)
	var mgoData []map[string]interface{}
	queryRS.All(&mgoData)
	return mgoData
}
