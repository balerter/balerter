package chart

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
	"io"
	"log"
)

const (
	defaultChartWidth  = 600
	defaultChartHeight = 200
)

type Chart struct {
	title  string
	series []Series
}

func New(title string) *Chart {
	ch := &Chart{
		title: title,
	}

	return ch
}

type Datum struct {
	Timestamp int64
	Value     float64
}

type Series struct {
	Title string
	Data  []Datum
}

func (ch *Chart) AddSeries(s Series) {
	ch.series = append(ch.series, s)
}

func (ch *Chart) Render(w io.Writer) error {

	p, err := plot.New()
	if err != nil {
		return fmt.Errorf("error create new plot, %w", err)
	}

	xticks := plot.TimeTicks{Format: "2006-01-02\n15:04"}

	p.Title.Text = ch.title
	//p.X.Label.Text = "X"
	//p.Y.Label.Text = "Y"
	p.X.Tick.Marker = xticks
	//p.Y.Min = 0
	//p.Y.Max = 1
	//p.Add(plotter.NewGrid())

	for seriesID, series := range ch.series {
		_ = seriesID
		//cnt := 0
		data := make(plotter.XYs, 0)
		for _, value := range series.Data {
			xy := plotter.XY{
				X: float64(value.Timestamp),
				Y: value.Value,
			}
			data = append(data, xy)
		}

		line, points, err := plotter.NewLinePoints(data)
		if err != nil {
			log.Panic(err)
		}
		//line.Color = color.RGBA{R: 0, G: 0, B: 180, A: 255}
		//points.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255}
		points.Shape = draw.CircleGlyph{}

		p.Add(line, points)
	}

	ww, err := p.WriterTo(defaultChartWidth, defaultChartHeight, "png")
	if err != nil {
		return err
	}

	_, err = ww.WriteTo(w)

	return err
}
