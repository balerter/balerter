package chart

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
	"io"
	"strconv"
	"strings"
)

const (
	defaultChartWidth  = 600
	defaultChartHeight = 200
)

var (
	defaultColors = map[string]color.RGBA{
		"blue":   color.RGBA{R: 0, G: 0, B: 255, A: 255},
		"red":    color.RGBA{R: 255, G: 0, B: 0, A: 255},
		"black":  color.RGBA{R: 0, G: 0, B: 0, A: 255},
		"green":  color.RGBA{R: 0, G: 255, B: 0, A: 255},
		"yellow": color.RGBA{R: 255, G: 255, B: 0, A: 255},
	}
)

func (ch *Chart) parseColor(s string) (color.RGBA, error) {
	var c color.RGBA

	if s == "" {
		return color.RGBA{A: 255}, nil
	}

	if c, ok := defaultColors[s]; ok {
		return c, nil
	}

	if !strings.HasPrefix(s, "#") {
		return c, fmt.Errorf("wrong color format")
	}

	s = s[1:]

	switch len(s) {
	case 6:
		return parseColor6(s)
	case 8:
		return parseColor8(s)
	}

	return c, fmt.Errorf("wrong color format")
}

func parseGroup(s string) (uint8, error) {
	if len(s) != 2 {
		return 0, fmt.Errorf("wrong group length")
	}

	r, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("wrong group format, %w", err)
	}

	return uint8(r), nil
}

func parseColor6(s string) (color.RGBA, error) {
	var c color.RGBA
	r, err := parseGroup(s[:2])
	if err != nil {
		return c, err
	}
	g, err := parseGroup(s[2:4])
	if err != nil {
		return c, err
	}
	b, err := parseGroup(s[4:])
	if err != nil {
		return c, err
	}

	c.R = r
	c.G = g
	c.B = b
	c.A = 255

	return c, nil
}

func parseColor8(s string) (color.RGBA, error) {
	c, err := parseColor6(s[:6])
	if err != nil {
		return c, err
	}

	a, err := parseGroup(s[6:])
	if err != nil {
		return c, err
	}

	c.A = a

	return c, nil
}

func (ch *Chart) Render(title string, data *Data, w io.Writer) error {

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
		for _, value := range series.Data {
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

		if series.LineColor == "" {
			series.LineColor = series.Color
		}
		if series.PointColor == "" {
			series.PointColor = series.Color
		}

		line.Color, err = ch.parseColor(series.LineColor)
		if err != nil {
			return fmt.Errorf("error parse line color, %w", err)
		}
		points.Color, err = ch.parseColor(series.PointColor)
		if err != nil {
			return fmt.Errorf("error parse point color, %w", err)
		}

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
