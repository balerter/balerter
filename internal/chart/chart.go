package chart

import (
	"math/rand"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
)

type Chart struct {
}

func New() *Chart {
	ch := &Chart{}

	return ch
}

func (ch *Chart) Render() {
	rand.Seed(int64(0))

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Plotutil example"
	p.X.Label.Text = " "
	p.Y.Label.Text = " "

	err = plotutil.AddLinePoints(p,
		"", randomPoints(15),
	)
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(600, 200, "points.png"); err != nil {
		panic(err)
	}
}

func randomPoints(n int) plotter.XYs {
	pts := make(plotter.XYs, n)
	for i := range pts {
		if i == 0 {
			pts[i].X = rand.Float64()
		} else {
			pts[i].X = pts[i-1].X + rand.Float64()
		}
		pts[i].Y = pts[i].X + 10*rand.Float64()
	}
	return pts
}
