package main

import (
    "image"
    "image/draw"
    "image/jpeg"
    "image/png"
	"os"
	"path/filepath"
    "fmt"
    "strings"
    "math"
	"strconv"

    "github.com/andlabs/ui"
    _ "github.com/andlabs/ui/winmanifest"
)

var mainwin *ui.Window
var entry_intput_dir *ui.Entry
var entry_output_file *ui.Entry
var entry_slice_w  *ui.Entry
var entry_slice_h  *ui.Entry
var entry_small_w  *ui.Entry
var entry_small_h  *ui.Entry
var bar_process *ui.ProgressBar
var btn_ok *ui.Button

func get_common_divisor(x int, y int) int {
    if x == y {
        return x
    }
    if x < y {
        x, y = y, x
    }
    if y == 0 {
        return x;
    }
    z := x % y
    return get_common_divisor(y, z)
}

func merge(path string, slice_width int, slice_height int, small_width int, small_height int, outputName string) error {
    defer func() {
        btn_ok.Enable()
    }()

    var image_paths []string
	err := filepath.Walk(path, func (subPath string, info os.FileInfo, err error) error {
		if err != nil {
            ui.MsgBox(mainwin, "error", err.Error())
			return err
		}
        if len(subPath) < 4 ||
        ! strings.EqualFold(subPath[len(subPath)-4:len(subPath)], ".png") {
            return nil
        }
        image_paths  = append(image_paths, subPath)
        return nil
    })

    image_slice_count := len(image_paths)
    fmt.Println(image_slice_count, slice_height, small_width, slice_width, small_height)
    h_count := int(math.Sqrt(float64(image_slice_count * slice_height * small_width) / float64(slice_width * small_height)))
    v_count := image_slice_count / h_count
    fmt.Println(h_count, v_count)

    output_width := slice_width * h_count
    output_height := slice_height * v_count
    output_image := image.NewRGBA(image.Rect(0, 0, output_width, output_height))
    x := 0
    y := 0
    for i, subPath := range image_paths {
		//fmt.Printf(subPath)
		imgb, _ := os.Open(subPath)
        defer imgb.Close()
        img, _ := png.Decode(imgb)
        draw.Draw(output_image, image.Rect(x, y, x + slice_width, y + slice_height), img, image.ZP, draw.Src)
        y += slice_height
        if y == output_height {
            x += slice_width
            y = 0
        }
        //fmt.Println(" success.")
        bar_process.SetValue(int(float64(i + 1)/float64(len(image_paths)) * 100))
    }
    
    output_file, err := os.Create(outputName)
    defer output_file.Close()
    if err != nil {
        fmt.Println(err)
        ui.MsgBox(mainwin, "error", err.Error())
        return err
	}
	jpeg.Encode(output_file, output_image, &jpeg.Options{jpeg.DefaultQuality})
	return nil
}

func output(b *ui.Button) {
    fmt.Println(btn_ok.Enabled(), btn_ok.Visible())
    if !btn_ok.Enabled() || !btn_ok.Visible() {
        return
    }
	slice_w,err1 := strconv.Atoi(entry_slice_w.Text())
	if err1 != nil {
		ui.MsgBox(mainwin, "参数错误", "宽高请输入数字")
		return
	}
	slice_h,err2 := strconv.Atoi(entry_slice_h.Text())
	if err2 != nil  {
		ui.MsgBox(mainwin, "参数错误", "宽高请输入数字")
		return
	}
	small_w,err3 := strconv.Atoi(entry_small_w.Text())
	if err3 != nil  {
		ui.MsgBox(mainwin, "参数错误", "宽高请输入数字")
		return
	}
	small_h,err4 := strconv.Atoi(entry_small_h.Text())
	if err4 != nil  {
		ui.MsgBox(mainwin, "参数错误", "宽高请输入数字")
		return
    }
    btn_ok.Disable()
	go merge(entry_intput_dir.Text(),
	slice_w,
	slice_h,
	small_w,
	small_h,
	entry_output_file.Text())
}

func setupUI() {
	mainwin = ui.NewWindow("合成图片", 480, 240, true)
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	vbox.Append(entryForm, false)

	entry_intput_dir = ui.NewEntry()
	entry_output_file = ui.NewEntry()

	entryForm.Append("导入路径", entry_intput_dir, false)
	entryForm.Append("导出文件", entry_output_file, false)

	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	vbox.Append(hbox, false)

    form_slice := ui.NewForm()
    form_slice.SetPadded(true)
    hbox.Append(form_slice, false)
    
    form_small := ui.NewForm()
    form_small.SetPadded(true)
	hbox.Append(form_small, false)

	entry_slice_w = ui.NewEntry()
	entry_slice_h = ui.NewEntry()
	entry_small_w = ui.NewEntry()
    entry_small_h = ui.NewEntry()
    bar_process = ui.NewProgressBar()
	btn_ok = ui.NewButton("OK")
    btn_ok.OnClicked(output)

	entry_slice_w.SetText("96")
	entry_slice_h.SetText("64")

	form_slice.Append("切图宽", entry_slice_w, false)
	form_slice.Append("切图高", entry_slice_h, false)
	form_small.Append("      小预览图宽", entry_small_w, false)
	form_small.Append("      小预览图高", entry_small_h, false)
	vbox.Append(bar_process, false)
	vbox.Append(btn_ok, false)

	mainwin.SetChild(vbox)
	mainwin.SetMargined(true)

	mainwin.Show()
}

func main() {
    ui.Main(setupUI)
    // enc := mahonia.NewEncoder("utf8")

    // app := app.New()
    // w := app.NewWindow("图片合成")
    // entry_input := widget.NewEntry()
    // entry_output := widget.NewEntry()
    // entry_slice_w := widget.NewEntry()
    // entry_slice_w.SetText("96")
    // entry_slice_h := widget.NewEntry()
    // entry_slice_h.SetText("64")
    // entry_small_w := widget.NewEntry()
    // entry_small_w.SetText("165")
    // entry_small_h := widget.NewEntry()
    // entry_small_h.SetText("180")
    // w.Resize(fyne.NewSize(800, 600))
	// w.SetContent(widget.NewVBox(
    //     widget.NewHBox(
    //         widget.NewLabel(enc.ConvertString("导入目录")),
    //         entry_input,
    //     ),
    //     widget.NewHBox(
    //         widget.NewLabel(enc.ConvertString("导出目录")),
    //         entry_output,
    //     ),
    //     widget.NewHBox(
    //         widget.NewLabel(enc.ConvertString("切图宽")),
    //         entry_slice_w,
    //         widget.NewLabel(enc.ConvertString("切图高")),
    //         entry_slice_h,
    //         widget.NewLabel(enc.ConvertString("小预览图宽")),
    //         entry_small_w,
    //         widget.NewLabel(enc.ConvertString("小预览图高")),
    //         entry_small_h,
    //     ),
	// 	widget.NewButton("OK", func() {
    //         merge("E:\\h5Proj3\\1\\", 96, 64, 165, 180, "1.png")
	// 	}),
	// ))
	// w.ShowAndRun()
}