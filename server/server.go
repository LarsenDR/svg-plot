package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"

	sf "github.com/larsendr/svgfunctions"
)

//Metadata structure is nested in to several structure

//ScreenData is one of base metadata structures
type ScreenData struct {
	Label  string
	Xorg   int
	Yorg   int
	Width  int
	Height int
}

//GraphData if the dat about the graph part of the
type GraphData struct {
	Label        string
	Space        string
	XaxisLabel   string
	XaxisUnitmax int
	XaxisUnitmin int
	YaxisLabel   string
	YaxisUnitmax int
	YaxisUnitmin int
	Grid         bool
	GridUnit     int
	GridColor    string
}

//MarginsData are a set of four side Margin structures
type MarginsData struct {
	Left  Margin
	Bott  Margin
	Right Margin
	Top   Margin
}

//Margin is a meta data structure about on side of the graph.
type Margin struct {
	Side          string
	Size          int
	AxisLine      bool
	Stroke        string
	StrokeWt      int
	Ticunit       int
	Ticsize       int
	Ticstroke     int
	Ticfontpx     int
	Ticfontoffset int
	Labelpx       int
	Labeltext     string
}

//Data is the graph general data structure
type Data struct {
	Testing                bool
	TestingBackgroundColor string
	TestingStrokeColor     string
	BackgroundColor        string
	StrokeColor            string
	FontFamily             string
	Screen                 ScreenData
	Graph                  GraphData
	Margs                  MarginsData
}

//ClientData is a structure to hold point to plot in graph space.
type ClientPoint struct {
	X int `json:"x"`
	Y int `json:"y"`
	R int `json:"r"`
}

type ClientPoints struct {
	DataVals []ClientPoint
}

//GetClientData is a function to read the datapoints.json into ClientData.
func GetClientData(filename string) (ClientPoints, error) {
	var cldt ClientPoints

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("%#v\n", err)
	}

	json.Unmarshal(file, &cldt)
	if err != nil {
		fmt.Printf("%#v\n", err)
	}
	log.Printf("%#v\n", cldt)
	return cldt, nil
}

//GetLayout is a function to read the graph.json settings
func GetLayout(filename string) (Data, error) {
	var data Data

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("%#v\n", err)
	}

	dec := json.NewDecoder(strings.NewReader(string(file)))
	err = dec.Decode(&data)
	if err != nil {
		fmt.Printf("%#v\n", err)
	}

	// log.Printf("%#v\n", data)

	return data, err
}

//TestHandler function
func (dt *Data) TestHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("New Graph\n")
	//var gx, gy int
	var width, height int
	var ScHeight int

	// Calculate screen aspect ratio
	var xunits, yunits float64
	var xscale, yscale float64
	var xgridsize, ygridsize int
	var gunit float64
	xunits = float64(dt.Graph.XaxisUnitmax) - float64(dt.Graph.XaxisUnitmin)
	yunits = float64(dt.Graph.YaxisUnitmax) - float64(dt.Graph.YaxisUnitmin)

	// If the graph units are the same  and the margins are the same the screens units should be the same
	ScHeight = int(float64(dt.Screen.Width)*(float64(yunits)/float64(xunits))) + (dt.Margs.Top.Size + dt.Margs.Bott.Size)

	width = dt.Screen.Width - (dt.Margs.Left.Size + dt.Margs.Right.Size)
	// to maintain the aspect ratio the height is adjusted to be proportional to the graph aspect ratio
	height = int(float64(width) * ((yunits) / (xunits)))

	xscale = float64(width) / float64(xunits)
	yscale = float64(height) / -(float64(yunits))

	gunit = float64(dt.Graph.GridUnit)
	xgridsize = int(xunits / gunit)
	ygridsize = int(yunits / gunit)

	log.Printf("Calculated xunits %f, yunits %f\n", xunits, yunits)

	log.Printf("Calculated xscale %f, yscale %f\n", xscale, yscale)

	log.Printf("Calculated aspect ratio %f\n", math.Abs(yscale)/xscale)

	log.Printf("Calculated Graph width %d, Graph height %d\n", width, height)

	log.Printf("Calculated X grid size %v, Y grid size %v\n", xgridsize, ygridsize)

	log.Printf("Calculated Screen Width %d, Screen Height %d\n", dt.Screen.Width, ScHeight)

	fmt.Fprintf(w, "%s\n", sf.XMLStart)
	fmt.Fprintf(w, "%s\n", sf.SVGStart(dt.Screen.Width, ScHeight))

	//Graph
	if dt.Graph.Grid {
		fmt.Fprintf(w, "%s\n", sf.SVGGrid(dt.Graph.GridColor, dt.Margs.Left.Size, dt.Margs.Top.Size, width, height, xgridsize, ygridsize))
	}

	//Top margin
	fmt.Fprintf(w, "<g transform=\"translate(%d %d)\">\n", dt.Margs.Left.Size, 0)
	if dt.Testing {
		//If testing this draws the margin box
		fmt.Fprintf(w, "%s\n", sf.SVGRect(dt.TestingBackgroundColor, dt.TestingStrokeColor, 0, 0, width, dt.Margs.Top.Size))
	}
	if dt.Margs.Top.AxisLine {
		fmt.Fprintf(w, "%s\n", sf.SVGLine(dt.Margs.Top.Stroke, dt.Margs.Top.StrokeWt, 0, dt.Margs.Top.Size, width, dt.Margs.Top.Size))
		for i := 0; i <= width; i += int(width / xgridsize) {
			fmt.Fprintf(w, "%s\n", sf.SVGLine(dt.Margs.Top.Stroke, dt.Margs.Top.StrokeWt, i, dt.Margs.Top.Size-dt.Margs.Top.Ticsize, i, dt.Margs.Top.Size))
			ytext := dt.Margs.Top.Size - (dt.Margs.Top.Ticsize + dt.Margs.Top.Ticfontoffset)
			fmt.Fprintf(w, "%s\n", sf.SVGTextTicLabel(i, ytext, "Top", dt.FontFamily, dt.Margs.Top.Ticfontpx, i))
		}
	}
	fmt.Fprintf(w, "<text x=\"%d\" y=\"%d\" font-family=\"%s\" font-size=\"%dpx\" text-anchor=\"middle\" alignment-basline=\"middle\" > %s </text>",
		width/2, dt.Margs.Top.Size/2, dt.FontFamily, dt.Margs.Top.Labelpx, dt.Margs.Top.Labeltext)
	fmt.Fprintf(w, " %s\n", `</g>`)

	//Bottom margin
	fmt.Fprintf(w, "<g transform=\"translate(%d %d)\">\n", dt.Margs.Left.Size, height+dt.Margs.Top.Size)
	if dt.Testing {
		//If testing this draws the margin box
		fmt.Fprintf(w, "%s\n", sf.SVGRect(dt.TestingBackgroundColor, dt.TestingStrokeColor, 0, 0, width, dt.Margs.Bott.Size))
	}
	if dt.Margs.Bott.AxisLine {
		fmt.Fprintf(w, "%s\n", sf.SVGLine(dt.Margs.Bott.Stroke, dt.Margs.Bott.StrokeWt, 0, 0, width, 0))
		// fmt.Fprintf(w, DrawMarginAxis(d.dt.Marg Bott, width))
		for i := 0; i <= width; i += int(width / xgridsize) {
			fmt.Fprintf(w, "%s\n", sf.SVGLine(dt.Margs.Bott.Stroke, dt.Margs.Bott.StrokeWt, i, 0, i, dt.Margs.Bott.Ticsize))
			fmt.Fprintf(w, "%s\n", sf.SVGTextTicLabel(i, dt.Margs.Bott.Ticsize+dt.Margs.Bott.Ticfontoffset, "Bott", dt.FontFamily, dt.Margs.Bott.Ticfontpx, i))
		}
	}
	fmt.Fprintf(w, "<text x=\"%d\" y=\"%d\" font-family=\"%s\" font-size=\"%dpx\" text-anchor=\"middle\" alignment-baseline=\"middle\" > %s </text>",
		width/2, dt.Margs.Bott.Size/2, dt.FontFamily, dt.Margs.Bott.Labelpx, dt.Margs.Bott.Labeltext)
	fmt.Fprintf(w, " %s\n", `</g>`)

	//Left margin
	fmt.Fprintf(w, "<g transform=\"translate(%d %d)\">\n", 0, dt.Margs.Top.Size)
	if dt.Testing {
		//If testing this draws the margin box
		fmt.Fprintf(w, "%s\n", sf.SVGRect(dt.TestingBackgroundColor, dt.TestingStrokeColor, 0, 0, dt.Margs.Left.Size, height))
	}
	if dt.Margs.Left.AxisLine {
		fmt.Fprintf(w, "%s\n", sf.SVGLine(dt.Margs.Left.Stroke, dt.Margs.Left.StrokeWt, dt.Margs.Left.Size, 0, dt.Margs.Left.Size, height))
		for i := 0; i <= height; i += int(height / ygridsize) {
			fmt.Fprintf(w, "%s\n", sf.SVGLine(dt.Margs.Left.Stroke, dt.Margs.Left.StrokeWt, dt.Margs.Left.Size, i, dt.Margs.Left.Size-dt.Margs.Left.Ticsize, i))
			xtext := dt.Margs.Left.Size - (dt.Margs.Left.Ticsize + dt.Margs.Left.Ticfontoffset)
			fmt.Fprintf(w, "%s\n", sf.SVGTextTicLabel(xtext, i, "Left", dt.FontFamily, dt.Margs.Left.Ticfontpx, height-i))
		}
	}
	fmt.Fprintf(w, "<text x=\"%d\" y=\"%d\" font-family=\"%s\" font-size=\"%dpx\" text-anchor=\"middle\" alignment-baseline=\"middle\" transform=\"rotate(270 %d %d)\"> %s </text>",
		dt.Margs.Bott.Size/2, height/2, dt.FontFamily, dt.Margs.Left.Labelpx, dt.Margs.Bott.Size/2, height/2, dt.Margs.Left.Labeltext)
	fmt.Fprintf(w, " %s\n", `</g>`)

	//Right margin
	fmt.Fprintf(w, "<g transform=\"translate(%d %d)\">\n", width+dt.Margs.Left.Size, dt.Margs.Top.Size)
	if dt.Testing {
		//If testing this draws the margin box
		fmt.Fprintf(w, "%s\n", sf.SVGRect(dt.TestingBackgroundColor, dt.TestingStrokeColor, 0, 0, dt.Margs.Right.Size, height))
	}
	if dt.Margs.Right.AxisLine {
		fmt.Fprintf(w, "%s\n", sf.SVGLine(dt.Margs.Right.Stroke, dt.Margs.Right.StrokeWt, 0, 0, 0, height))
		for i := 0; i <= height; i += int(height / ygridsize) {
			fmt.Fprintf(w, "%s\n", sf.SVGLine(dt.Margs.Right.Stroke, dt.Margs.Right.StrokeWt, 0, i, dt.Margs.Right.Ticsize, i))
			xtext := dt.Margs.Right.Ticsize + dt.Margs.Right.Ticfontoffset
			fmt.Fprintf(w, "%s\n", sf.SVGTextTicLabel(xtext, i, "Right", dt.FontFamily, dt.Margs.Right.Ticfontpx, height-i))
		}
	}

	fmt.Fprintf(w, "<text x=\"%d\" y=\"%d\" font-family=\"%s\" font-size=\"%dpx\" text-anchor=\"middle\" alignment-basline=\"middle\" transform=\"rotate(90 %d %d)\"> %s </text>",
		dt.Margs.Right.Size/2, height/2, dt.FontFamily, dt.Margs.Right.Labelpx, dt.Margs.Right.Size/2, height/2, dt.Margs.Right.Labeltext)
	fmt.Fprintf(w, " %s\n", `</g>`)

	log.Printf("Scale (%f %f )", xscale, yscale)

	//fmt.Printf("Scale output %v\n", ScaleMathToGraph(200, 0, 660, 0, 1000))

	// plot points
	// fmt.Fprintf(w, "<g transform=\"translate(%d %d) scale(%.3f %.3f)\">\n", d.Margin Left.Size, height+d.Margin Top.Size, xscale, yscale)
	// fmt.Fprintf(w, "%s\n", sf.SVGPoint("#ff0000", "#000000", 1, 10, 10, 2))

	// fmt.Fprintf(w, " %s\n", `</g>`)

	fmt.Fprintf(w, "%s\n", sf.SVGEnd)
}
