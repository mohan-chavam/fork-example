package main

import (
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/kuangcp/gobase/pkg/cuibase"
	"github.com/kuangcp/gobase/toolbox/keylogger/app/store"
	"log"
	"time"
	"unsafe"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/kuangcp/logger"
)

var (
	DashboardMsMode bool // 以 ms 格式刷新时间
	DashboardMs     = DefaultRefreshMs
)

const (
	width            = 84
	height           = 10
	DefaultRefreshMs = 57
	appId            = "com.github.kuangcp.keylogger"
)

var (
	app           *gtk.Application
	win           *gtk.Window
	kpmLabel      *gtk.Label
	refreshPeriod = time.Millisecond * 400
)

var (
	help bool

	// redis
	host string
	port string
	pwd  string
	db   int

	option redis.Options
)

var (
	buildVersion string
)

var info = cuibase.HelpInfo{
	Description:   "Record key input, show rank",
	Version:       "1.1.0",
	BuildVersion:  buildVersion,
	SingleFlagLen: -5,
	DoubleFlagLen: 0,
	ValueLen:      -6,
	Flags: []cuibase.ParamVO{
		{Short: "-h", BoolVar: &help, Comment: "help info"},
		{Short: "-m", BoolVar: &DashboardMsMode, Comment: "show time with ms"},
	},
	Options: []cuibase.ParamVO{
		{Short: "-host", Value: "host", Comment: "redis host"},
		{Short: "-port", Value: "port", Comment: "redis port"},
		{Short: "-pwd", Value: "pwd", Comment: "redis password"},
		{Short: "-db", Value: "db", Comment: "redis db"},
		{Short: "-ms", Value: "ms", Comment: "gui refresh ms"},
	},
}

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "")
	flag.StringVar(&port, "port", "6667", "")
	flag.StringVar(&pwd, "pwd", "", "")
	flag.IntVar(&db, "db", 5, "")
	flag.IntVar(&DashboardMs, "ms", DashboardMs, "")
}

func main() {
	info.Parse()

	if help {
		info.PrintHelp()
		return
	}

	option = redis.Options{Addr: host + ":" + port, Password: pwd, DB: db}
	store.InitConnection(option)

	if DashboardMsMode || (!DashboardMsMode && DashboardMs != DefaultRefreshMs) {
		refreshPeriod = time.Millisecond * time.Duration(DashboardMs)
	}
	logger.Info("refresh:", refreshPeriod)

	gtk.Init(nil)
	app, _ = gtk.ApplicationNew(appId, glib.APPLICATION_FLAGS_NONE)
	app.Connect("activate", createWindow)
	app.Run(nil)
}

func createWindow() {
	win, _ = gtk.WindowNew(gtk.WINDOW_POPUP)
	win.SetDefaultSize(width, height)
	win.SetPosition(gtk.WIN_POS_MOUSE)
	gridWidget := createLabelWidget()
	win.Add(gridWidget)
	win.Connect("destroy", gtk.MainQuit)
	bindMouseActionForWindow()

	app.AddWindow(win)
	win.ShowAll()

	// 启动后的计算并刷新缓存
	go timeoutRefresh(refreshPeriod)
}

// https://developer.gnome.org/gtk4/unstable/GtkLabel.html
func createLabelWidget() *gtk.Widget {
	grid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create grid:", err)
	}
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)

	kpmLabel, err = gtk.LabelNew("")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	grid.Add(kpmLabel)
	kpmLabel.SetMarkup(latestLabelStr(time.Now()))
	kpmLabel.SetHExpand(true)
	kpmLabel.SetVExpand(true)

	return &grid.Container.Widget
}

func bindMouseActionForWindow() {
	// 鼠标最后点击坐标
	var x, y int
	win.SetEvents(int(gdk.BUTTON_PRESS_MASK | gdk.BUTTON1_MOTION_MASK))

	//鼠标按下事件处理
	win.Connect("button-press-event", func(widget *gtk.Window, ctx *gdk.Event) {
		//获取鼠键按下属性结构体变量，系统内部的变量，不是用户传参变量
		event := *(*gdk.EventButton)(unsafe.Pointer(&ctx))
		if event.Button() == 1 { //左键
			x, y = int(event.X()), int(event.Y()) //保存点击的起点坐标
		} else if event.Button() == 3 { //右键
			//右键，关闭窗口
			app.Quit()
		}
	})

	// 鼠标移动事件处理
	// 注意： 不出现偏差的前提是 应用 内不出现其他点击事件和交互
	win.Connect("motion-notify-event", func(widget *gtk.Window, ctx *gdk.Event) {
		//获取鼠标移动属性结构体变量，系统内部的变量，不是用户传参变量
		event := *(*gdk.EventButton)(unsafe.Pointer(&ctx))
		win.Move(int(event.XRoot())-x, int(event.YRoot())-y)
	})
}

func latestLabelStr(now time.Time) string {
	conn := store.GetConnection()
	today := now.Format(store.DateFormat)

	tempValue, err := conn.Get(store.GetTodayTempKPMKeyByString(today)).Result()
	if err != nil {
		tempValue = "0"
	}
	maxValue, err := conn.Get(store.GetTodayMaxKPMKeyByString(today)).Result()
	if err != nil {
		maxValue = "0"
	}
	total := conn.ZScore(store.TotalCount, today).Val()

	var timeFmt = ""
	if DashboardMsMode {
		timeFmt = store.MsTimeFormat
	} else {
		timeFmt = store.TimeFormat
	}

	// style https://blog.csdn.net/bitscro/article/details/3874616
	return "<span font_family='Cascadia Mono PL' font_desc='10'>" +
		"<span foreground='#00FFF6'>" + fmt.Sprintf("%11s", now.Format(timeFmt)) + "</span>\n" +
		"<span foreground='#5AFF00'>" + fmt.Sprintf("%3s", tempValue) + "</span> " +
		"<span foreground='gray'>" + fmt.Sprintf("%3s", maxValue) + "</span> " +
		"<span foreground='white'>" + fmt.Sprintf("%-6d", int(total)) + "</span></span>"
}

// 从缓存中更新窗口内面板
func timeoutRefresh(period time.Duration) {
	// 返回 true 才能一直执行
	glib.TimeoutAdd(uint(period.Milliseconds()), func() bool {
		str := latestLabelStr(time.Now())
		kpmLabel.SetMarkup(str)
		return true
	})
}
