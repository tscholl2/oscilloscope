// +build ignore

package main

import (
	"time"

	ui "github.com/gizak/termui"
)

// returns a pointer to an array of bytes from
// the data source. It should match the length.
// this will be called every 1/4 second so it should
// take less time than that to fill up.
func dataSource() *[]byte {
	b := make([]byte, 111)
	//rand.Read(b)
	for i := 0; i < len(b); i++ {
		b[i] = byte(time.Now().Nanosecond())
	}
	return &b
}

// samples the data from source.
func sampleData(length int) *[]float64 {
	b := *dataSource()
	a := make([]float64, length)
	for i := range a {
		a[i] = float64(int(b[i*len(b)/len(a)]))
	}
	return &a
}

// shifts the new data into the left side
// so shift([1,2,3,4],[5,6]) = [3,4,5,6]
func shift(arr, stuff *[]float64) *[]float64 {
	alen := len(*arr)
	slen := len(*stuff)
	if alen < slen {
		panic("Cant shift more than it can hold")
	}
	for i := 0; i < alen-slen; i++ {
		(*arr)[i] = (*arr)[slen+i]
	}
	for i := alen - slen; i < alen; i++ {
		(*arr)[i] = (*stuff)[i-alen+slen]
	}
	return arr
}

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	ui.UseTheme("helloworld")

	sampleSize := 64

	lc := ui.NewLineChart()
	lc.Border.Label = "Oscilloscope"
	lc.Data = make([]float64, sampleSize)
	lc.Height = 24
	lc.AxesColor = ui.ColorWhite
	lc.LineColor = ui.ColorYellow | ui.AttrBold
	lc.Mode = "dot"

	// build layout
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(12, 0, lc)),
	)

	// calculate layout
	ui.Body.Align()

	redraw := make(chan bool)

	update := func() {
		for {
			lc.Data = *shift(&lc.Data, sampleData(sampleSize/8))
			time.Sleep(time.Second / 4)
			redraw <- true
		}
	}

	evt := ui.EventCh()

	ui.Render(ui.Body)
	go update()

	for {
		select {
		case e := <-evt:
			if e.Type == ui.EventKey && e.Ch == 'q' {
				return
			}
			if e.Type == ui.EventResize {
				ui.Body.Width = ui.TermWidth()
				ui.Body.Align()
				go func() { redraw <- true }()
			}
		case <-redraw:
			ui.Render(ui.Body)
		}
	}
}
