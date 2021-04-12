package html

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

var htmlTemplate = template.Must(template.New("html").Funcs(map[string]interface{}{
	"formatFloat": func(f float64, decimals int) string {
		return strings.Replace(strconv.FormatFloat(f, 'f', decimals, 64), ".", ",", 1)
	},
	"toJSON": func(v interface{}) string {
		m, _ := json.Marshal(v)
		return string(m)
	},
}).Parse(tmplHTML))

const tmplHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <title>Results</title>
	<script src="https://cdn.tutorialjinni.com/sigma.js/1.2.1/sigma.min.js"></script>
	<script src="https://cdn.tutorialjinni.com/sigma.js/1.2.1/plugins/sigma.exporters.svg.min.js"></script>
	<script src="https://cdn.tutorialjinni.com/sigma.js/1.2.1/plugins/sigma.layout.forceAtlas2.min.js"></script>
	<script src="https://cdn.tutorialjinni.com/sigma.js/1.2.1/plugins/sigma.layout.noverlap.min.js"></script>
	<script src="https://cdn.tutorialjinni.com/sigma.js/1.2.1/plugins/sigma.neo4j.cypher.min.js"></script>
	<script src="https://cdn.tutorialjinni.com/sigma.js/1.2.1/plugins/sigma.parsers.json.min.js"></script>
	<script src="https://cdn.tutorialjinni.com/sigma.js/1.2.1/plugins/sigma.pathfinding.astar.min.js"></script>
	<script src="https://cdn.tutorialjinni.com/sigma.js/1.2.1/plugins/sigma.plugins.animate.min.js"></script>
	<script src="https://cdn.tutorialjinni.com/sigma.js/1.2.1/plugins/sigma.plugins.dragNodes.min.js"></script>
	<script src="https://cdn.tutorialjinni.com/sigma.js/1.2.1/plugins/sigma.plugins.filter.min.js"></script>
	<script src="https://cdn.tutorialjinni.com/sigma.js/1.2.1/plugins/sigma.plugins.filter.min.js"></script>
	<script src="https://cdn.tutorialjinni.com/sigma.js/1.2.1/plugins/sigma.plugins.neighborhoods.min.js"></script>
	<script src="https://cdn.tutorialjinni.com/sigma.js/1.2.1/plugins/sigma.plugins.relativeSize.min.js"></script>
	<script src="https://cdn.tutorialjinni.com/sigma.js/1.2.1/plugins/sigma.renderers.customEdgeShapes.min.js"></script>
	<script src="https://cdn.tutorialjinni.com/sigma.js/1.2.1/plugins/sigma.renderers.customShapes.min.js"></script>
	<script src="https://cdn.tutorialjinni.com/sigma.js/1.2.1/plugins/sigma.renderers.edgeDots.min.js"></script>
	<script src="https://cdn.tutorialjinni.com/sigma.js/1.2.1/plugins/sigma.renderers.edgeLabels.min.js"></script>
	<script src="https://cdn.tutorialjinni.com/sigma.js/1.2.1/plugins/sigma.renderers.parallelEdges.min.js"></script>
	<script src="https://cdn.tutorialjinni.com/sigma.js/1.2.1/plugins/sigma.renderers.snapshot.min.js"></script>
	<script src="https://cdn.tutorialjinni.com/sigma.js/1.2.1/plugins/sigma.statistics.HITS.min.js"></script>
</head>
	<body>
		<style>
			.table {
			  font-family: Arial, Helvetica, sans-serif;
			  border-collapse: collapse;
			  width: 100%;
			}
			
			.table td, .table th {
			  border: 1px solid #ddd;
			  padding: 8px;
			}
			
			.table tr:nth-child(even){background-color: #f2f2f2;}
			
			.table tr:hover {background-color: #ddd;}
			
			.table th {
			  padding-top: 12px;
			  padding-bottom: 12px;
			  text-align: left;
			  background-color: #4CAF50;
			  color: white;
			}
			</style>
		<div class="left">
			<div>
				<table class="table">
					<thead>
						<tr>
							<th>#</td>
							<th>Implementation</td>
						</tr>
					</thead>
					<tbody>
						{{ range $index, $element := .Implementations }}
						<tr>
							<td>{{$index}}</td>
							<td>
								{{ range $element }}
									{{ . }} →
								{{ end }}
								...
							</td>
						</tr>
						{{ end}}
					</tbody>
				</table>
			</div>
			{{ range $index, $info :=  .Diffs }}
				<div class="table">
					<h3>Step #{{ $index }}</h3>
					<table class="fl-table">
						<thead>
							<tr>
								<th>#</td>
								<th>theoretical</td>
								<th>empiric</td>
								<th>diff</td>
							</tr>
						</thead>
						<tbody>
						{{ range $index2, $i := $info }}
							<tr>
								<td>{{ $index2 }}</td>
								<td>{{ formatFloat $i.T 6 }}</td>
								<td>{{ formatFloat $i.E 6 }}</td>
								<td>{{ formatFloat $i.D 6 }}</td>
							</tr>
						{{ end }}
						</tbody>
					</table>
				</div>
			{{ end }}
		</div>
		<div class="right">
			<style>
				body {
				  display: flex;
				  width: 100%;
				  height: 100%;
				  margin: 0;
				  padding: 0;
				}
				body > * {
				  width: 50%;
				}
				.left {
				  border-right: 1px solid black;
				}
				.right {
					min-height: 500px;
				}
				.container {
				  	background-color: #fafafa;
				  	width: 100%;
				  	height: 100%;
					min-height: 500px;
				  	max-height: 100vh;
				}
				#graph-container{
					width: 100%;
					height: 100%;
					min-height: 500px;
				}
			</style>
			<div class="container">
			  <div id="graph-container"></div>
			</div>
		</div>
		<script>
			s = new sigma({
			  graph: JSON.parse({{toJSON .Graph}}),
			  renderer: {
				container: document.getElementById('graph-container'),
				type: 'canvas'
			  },
			  settings: {
				edgeLabelSize: 'proportional',
				enableEdgeHovering: true,
				edgeHoverSizeRatio: 2
			  }
			});
			const dragListener = sigma.plugins.dragNodes(s, s.renderers[0]);
		</script>
	</body>
</html>
	`

func Output(logger *log.Logger, d PageData) error {
	dir, err := ioutil.TempDir("", "result")
	if err != nil {
		return err
	}

	out, err := os.Create(filepath.Join(dir, "result.html"))
	if err != nil {
		return err
	}

	err = htmlTemplate.Execute(out, d)
	if err == nil {
		err = out.Close()
	}

	if err != nil {
		return err
	}

	if !startBrowser("file://" + out.Name()) {
		logger.Printf("HTML output written to %s\n", out.Name())
	}

	return nil
}

func startBrowser(url string) bool {
	// try to start the browser
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}
