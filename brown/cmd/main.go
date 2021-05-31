package main

import (
	"fmt"
	"github.com/Sinu5oid/generators/brown/rmd"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"math"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	n := 6
	sigma := float64(1)
	H := 0.75
	N := 100000
	NToDisplay := int(math.Min(float64(N), 10))
	engine := rmd.NewEngine(n, sigma, H)

	impls := make([][]float64, 0, N)
	for i := 0; i < N; i++ {
		impls = append(impls, engine.NextImpl())
	}
	plotImplementations(impls, NToDisplay, n)

	meansE := engine.MeanE(impls)
	meansT := engine.MeanT(impls)
	dispE := engine.DispE(impls, meansE)
	dispT := engine.DispT(impls)
	corellE := engine.CorellE(impls, dispE)
	corellT := engine.CorellT(impls, H)

	plotMeans(meansE, meansT, n)
	plotDispersion(dispE, n, dispT)
	for i := 0; i < len(corellE); i++ {
		plotCorells(corellE[i], corellT[i], n, fmt.Sprintf("corell_%d", i))
	}
}

func plotDispersion(dispE []float64, n int, dispT []float64) {
	p, err := plot.New()
	if err != nil {
		fmt.Println("failed to create plot", err)
		return
	}
	p.Title.Text = "Dispersion"
	p.X.Label.Text = "tk"
	p.Y.Label.Text = "d(tk)"

	p.Legend.Top = true

	liner, _, err := plotter.NewLinePoints(getXYs(dispE, n))
	if err != nil {
		fmt.Println("failed to create liner:", err)
	}

	liner.Color = plotutil.Color(0)
	liner.Dashes = plotutil.Dashes(0)

	p.Add(liner)
	p.Legend.Add("disp_e", liner)

	liner, _, err = plotter.NewLinePoints(getXYs(dispT, n))
	if err != nil {
		fmt.Println("failed to create liner:", err)
	}

	liner.Color = plotutil.Color(1)
	liner.Dashes = plotutil.Dashes(1)

	p.Add(liner)
	p.Legend.Add("disp_t", liner)

	filename := fmt.Sprintf("disp_output.png")

	if err := p.Save(10*vg.Inch, 5*vg.Inch, filename); err != nil {
		fmt.Println("failed to save file ", filename, err)
	}
}

func plotMeans(meansE, meansT []float64, n int) {
	p, err := plot.New()
	if err != nil {
		fmt.Println("failed to create plot", err)
		return
	}
	p.Title.Text = "Means"
	p.X.Label.Text = "tk"
	p.Y.Label.Text = "m(tk)"

	p.Legend.Top = true

	liner, _, err := plotter.NewLinePoints(getXYs(meansE, n))
	if err != nil {
		fmt.Println("failed to create liner:", err)
	}

	liner.Color = plotutil.Color(0)
	liner.Dashes = plotutil.Dashes(0)

	p.Add(liner)
	p.Legend.Add("mean_e", liner)

	liner, _, err = plotter.NewLinePoints(getXYs(meansT, n))
	if err != nil {
		fmt.Println("failed to create liner:", err)
	}

	liner.Color = plotutil.Color(1)
	liner.Dashes = plotutil.Dashes(1)

	p.Add(liner)
	p.Legend.Add("mean_t", liner)

	filename := fmt.Sprintf("mean_output.png")

	if err := p.Save(10*vg.Inch, 5*vg.Inch, filename); err != nil {
		fmt.Println("failed to save file ", filename, err)
	}
}

func plotCorells(correllE, corellT []float64, n int, name string) {
	p, err := plot.New()
	if err != nil {
		fmt.Println("failed to create plot", err)
		return
	}
	p.Title.Text = "Correlations"
	p.X.Label.Text = "tk"
	p.Y.Label.Text = "r(tk, tk + j)"

	p.Legend.Top = true

	liner, _, err := plotter.NewLinePoints(getXYs(correllE, n))
	if err != nil {
		fmt.Println("failed to create liner:", err)
	}

	liner.Color = plotutil.Color(0)
	liner.Dashes = plotutil.Dashes(0)

	p.Add(liner)
	p.Legend.Add("corell_e", liner)

	liner, _, err = plotter.NewLinePoints(getXYs(corellT, n))
	if err != nil {
		fmt.Println("failed to create liner:", err)
	}

	liner.Color = plotutil.Color(1)
	liner.Dashes = plotutil.Dashes(1)

	p.Add(liner)
	p.Legend.Add("corell_t", liner)

	filename := fmt.Sprintf("%s_output.png", name)

	if err := p.Save(10*vg.Inch, 5*vg.Inch, filename); err != nil {
		fmt.Println("failed to save file ", filename, err)
	}
}

func plotImplementations(impls [][]float64, NToDisplay int, n int) {
	p, err := plot.New()
	if err != nil {
		fmt.Println("failed to create plot", err)
		return
	}
	p.Title.Text = "Implementations"
	p.X.Label.Text = "tk"
	p.Y.Label.Text = "x(tk)"

	p.Legend.Top = true

	pts := getXYsGroup(impls[:NToDisplay], n)

	for i := 0; i < len(pts); i++ {
		liner, _, err := plotter.NewLinePoints(pts[i])
		if err != nil {
			fmt.Println("Failed to create liner:", err)
			continue
		}

		liner.Color = plotutil.Color(i)
		liner.Dashes = plotutil.Dashes(i)

		p.Add(liner)
		p.Legend.Add(fmt.Sprintf("impl#%d", i), liner)
	}

	filename := fmt.Sprintf("output.png")

	if err := p.Save(10*vg.Inch, 5*vg.Inch, filename); err != nil {
		fmt.Println("failed to save file ", filename, err)
		return
	}
}

func getXYsGroup(implementations [][]float64, n int) []plotter.XYer {
	xys := make([]plotter.XYer, 0, len(implementations))
	for i := 0; i < len(implementations); i++ {
		buffer := make(plotter.XYs, 0, len(implementations[i]))

		for k := 0; k < len(implementations[i]); k++ {
			buffer = append(buffer, plotter.XY{
				X: float64(k) / math.Pow(2, float64(n)),
				Y: implementations[i][k],
			})
		}

		xys = append(xys, buffer)
	}

	return xys
}

func getXYs(slice []float64, n int) plotter.XYs {
	xys := make(plotter.XYs, 0, len(slice))
	for k := 0; k < len(slice); k++ {
		xys = append(xys, plotter.XY{
			X: float64(k) / math.Pow(2, float64(n)),
			Y: slice[k],
		})
	}

	return xys
}
