package main

import (
	"github.com/Sinu5oid/generators/stochastic"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"

	"fmt"
	"image/color"
	"math/rand"
)

type Task struct {
	h            float64
	name         string
	index        int
	xLabel       string
	yLabel       string
	ideal        *[]*float64
	observed     *[]*float64
	yMin         *float64
	yMax         *float64
	outputFolder string
}

func CreateTask(
	h float64,
	name string,
	index int,
	xLabel string,
	yLabel string,
	ideal *[]*float64,
	observed *[]*float64,
	yMin *float64,
	yMax *float64,
	outputFolder string,
) *Task {
	return &Task{
		h:            h,
		name:         name,
		index:        index,
		xLabel:       xLabel,
		yLabel:       yLabel,
		ideal:        ideal,
		observed:     observed,
		yMin:         yMin,
		yMax:         yMax,
		outputFolder: outputFolder,
	}
}

func (t *Task) BuildPlot() error {
	p, err := plot.New()
	if err != nil {
		fmt.Println("failed to create plot", err)
		return err
	}
	p.Title.Text = t.name
	p.X.Label.Text = t.xLabel
	p.Y.Label.Text = t.yLabel

	if t.yMin != nil {
		p.Y.Min = *t.yMin
	}
	if t.yMin != nil {
		p.Y.Max = *t.yMax
	}

	p.Legend.Top = true

	err = plotutil.AddLines(p, "ideal", getXYs(t.ideal, t.h))
	if err != nil {
		fmt.Println("failed to add ideal lines to plot", err)
		return err
	}

	scatter, err := plotter.NewScatter(getXYs(t.observed, t.h))
	if err != nil {
		fmt.Println("failed to add ideal lines to plot", err)
		return err
	}
	scatter.GlyphStyle.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	scatter.GlyphStyle.Radius = vg.Points(1)

	p.Add(scatter)
	p.Legend.Add("observed", scatter)

	filename := fmt.Sprintf("%s/%s-%d.png", t.outputFolder, t.name, t.index)

	if err := p.Save(10*vg.Inch, 5*vg.Inch, filename); err != nil {
		fmt.Println("failed to save file ", filename, err)
		return err
	}

	return nil
}

func getXYs(src *[]*float64, h float64) *plotter.XYs {
	pts := make(plotter.XYs, len(*src), len(*src))

	for i := 0; i < len(*src); i += 1 {
		pts[i].X = stochastic.GetT(i, h)
		pts[i].Y = *(*src)[i]
	}

	return &pts
}

type SimpleTask struct {
	h            float64
	name         string
	xLabel       string
	yLabel       string
	sources      *[]*[]*float64
	outputFolder string
}

func CreateSimpleTask(
	h float64,
	name string,
	xLabel string,
	yLabel string,
	sources *[][]float64,
	outputFolder string,
) *SimpleTask {
	src := make([]*[]*float64, 0, len(*sources))
	for i := range *sources {
		subSrc := make([]*float64, 0, len((*sources)[i]))
		for j := range (*sources)[i] {
			subSrc = append(subSrc, &(*sources)[i][j])
		}
		src = append(src, &subSrc)
	}

	return &SimpleTask{
		h:            h,
		name:         name,
		xLabel:       xLabel,
		yLabel:       yLabel,
		sources:      &src,
		outputFolder: outputFolder,
	}
}

func (t *SimpleTask) BuildPlot() error {
	p, err := plot.New()
	if err != nil {
		fmt.Println("failed to create plot", err)
		return err
	}
	p.Title.Text = t.name
	p.X.Label.Text = t.xLabel
	p.Y.Label.Text = t.yLabel

	p.Legend.Top = true

	for i, source := range *t.sources {
		liner, scatter, err := plotter.NewLinePoints(getXYs(source, t.h))
		if err != nil {
			fmt.Println("failed to add ideal lines to plot", err)
			return err
		}

		c := color.RGBA{R: uint8(rand.Intn(255)), G: uint8(rand.Intn(255)), B: uint8(rand.Intn(255)), A: 255}
		liner.Color = c
		scatter.GlyphStyle.Color = c
		scatter.GlyphStyle.Radius = vg.Points(1)
		p.Legend.Add(fmt.Sprintf("x%d(t)", i), liner)

		p.Add(liner)
		p.Add(scatter)
	}

	filename := fmt.Sprintf("%s/%s.png", t.outputFolder, t.name)

	if err := p.Save(10*vg.Inch, 5*vg.Inch, filename); err != nil {
		fmt.Println("failed to save file ", filename, err)
		return err
	}

	return nil
}

type GenericTask interface {
	BuildPlot() error
}
