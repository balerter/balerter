package chart

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
	"io"
)

const (
	defaultChartWidth  = 600
	defaultChartHeight = 200
)

func (ch *Chart) _render(title string, data *Data, w io.Writer) error {

	p, err := plot.New()
	if err != nil {
		return fmt.Errorf("error create new plot, %w", err)
	}

	xticks := plot.TimeTicks{Format: "2006-01-02\n15:04"}

	p.Title.Text = title
	//p.X.Label.Text = "X"
	//p.Y.Label.Text = "Y"
	p.X.Tick.Marker = xticks
	//p.Y.Min = 0
	//p.Y.Max = 1
	//p.Add(plotter.NewGrid())

	for _, series := range data.Series {
		data := make(plotter.XYs, 0)
		for _, value := range series.Values {
			xy := plotter.XY{
				X: value.Timestamp,
				Y: value.Value,
			}
			data = append(data, xy)
		}

		line, points, err := plotter.NewLinePoints(data)
		if err != nil {
			return fmt.Errorf("error create line points, %w", err)
		}
		//line.Color = color.RGBA{R: 0, G: 0, B: 180, A: 255}
		//points.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255}
		points.Shape = draw.CircleGlyph{}

		p.Add(line, points)
	}

	ww, err := p.WriterTo(defaultChartWidth, defaultChartHeight, "png")
	if err != nil {
		return fmt.Errorf("error render data, %w", err)
	}

	_, err = ww.WriteTo(w)
	if err != nil {
		return fmt.Errorf("error write rendered data, %w", err)
	}

	return nil
}
