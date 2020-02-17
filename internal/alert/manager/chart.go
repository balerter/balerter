package manager

import (
	"bytes"
	"github.com/balerter/balerter/internal/chart"
	"io"
	"io/ioutil"
)

func (m *Manager) makeChart(opts *optionsChart) (string, error) {
	buf := bytes.NewBuffer([]byte{})

	err := m.renderImage(buf, opts)
	if err != nil {
		return "", err
	}

	// todo: change to upload
	err = ioutil.WriteFile("chart.png", buf.Bytes(), 0644)
	if err != nil {
		return "", err
	}

	return "", nil
}

func (m *Manager) renderImage(w io.Writer, opts *optionsChart) error {
	ch := chart.New(opts.Title)

	for _, series := range opts.Series {
		chartSeries := chart.Series{}
		for _, datum := range series.Data {
			chartSeries.Data = append(chartSeries.Data, chart.Datum{
				Timestamp: datum.Timestamp,
				Value:     datum.Value,
			})
		}
		ch.AddSeries(chartSeries)
	}

	return ch.Render(w)
}
