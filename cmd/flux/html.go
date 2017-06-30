package main

import (
	"html/template"
)

var (
	tmpl = template.Must(template.New("index").Parse(html))
)

const html = `
<!DOCTYPE html>
<html lang="en">
  	<head>
    	<meta charset="utf-8">
    	<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Flux - Dashboard</title>
		<script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/d3/3.5.17/d3.min.js"></script>
		<script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/c3/0.4.13/c3.min.js"></script>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/c3/0.4.13/c3.min.css"/>
		<style>
			body {
				background: #f8f8f8;
				margin: 0;
			}
			.container {
				width: 1200px;
				margin: 2em auto 4em auto;
			}
			.chart {
				background: #fff;
				width: 550px;
				float: left;
				padding: 15px;
				margin: 10px;
				box-shadow: 0 0 5px #888;
			}
			.overview {
				
			}
			body > header {
				color: #fff;
				background: #4958B8;
				padding: 1em 3em;
			}
			body > header h1 {
				font-weight: normal;
				margin: 0;
			}
		</style>
		<script>
		  	var __DATA__ = {{ . }};
		</script>
	</head>
  	<body>
	  	<header>
		  <h1>Flux</h1>
		</header>
		<div class="container">
			<div class="overview" id="overview-chart">
			</div>
			<div class="chart">
				<h3>Numar de inserari</h3>
				<div id="inserts-chart">
				</div>
			</div>
			<div class="chart">
				<h3>Numar de citiri</h3>
				<div id="reads-chart">
				</div>
			</div>
			<div class="chart">
				<h3>Memorie</h3>
				<div id="memory-chart">
				</div>
			</div>
			<div class="chart">
				<h3>Numar de chei</h3>
				<div id="keys-chart">
				</div>
			</div>
			<!--
			<div class="chart">
				<h3>Numar de stergeri</h3>
				<div id="deletions-chart">
				</div>
			</div>
			-->
		</div>
	  	<script>
			var chartTypes = [
				{ type: 'memory', label: 'Memory [MB]', getter: function(v) { return v / (1024*1024); } }, 
				{ type: 'keys', label: 'Keys' }, 
				{ type: 'inserts', label: 'Inserts' }, 
				{ type: 'reads', label: 'reads' }
			];
		
			function get_columns(data, type, getter) {
				var output = [];

				Object.keys(data.nodes).forEach(function(node) {
					var col = [node];
					Object.keys(data.metrics).forEach(function(date) {
						var val = data.metrics[date][node] ? data.metrics[date][node][type] : 0;
						col.push(getter? getter(val) : val)
					});
					output.push(col);
				});
				return output;
			}

		  	function draw(type, label, columns) {
				c3.generate({
					bindto: '#' + type + '-chart',
					data: {
						x: 'x',
						columns: columns,
					},
					point: {
						show: false
					},
					axis: {
						x: {
							type: 'timeseries',
							tick: {
								format: '%M:%S',
								count: 10,
								outer: true
							}
						},
						y: {
							label: {
								text: label,
								position: 'outer-middle'
							}
						}
					}
				});
			};
			
			function overview() {
				var sum = 0;
				var date = Object.keys(__DATA__.metrics).reduce(function(max, cur){ return max > cur ? max : cur });
				var keys = Object.keys(__DATA__.metrics[date]).map(function(key) {
					var k = __DATA__.nodes[key] + "(" + key + ")";
					var v = __DATA__.metrics[date][key].keys;

					sum = sum + v;
					
					return [k, v]
				});
				
				c3.generate({
					bindto: '#overview-chart',
					data: {
						columns: keys,
						type: 'donut'
					},
					donut: {
						title: sum
					}
				});
			};

			(function () {
				overview();
				var x = ['x'].concat(Object.keys(__DATA__.metrics).map(k => new Date(k)));
				
				for (var i = 0; i < chartTypes.length; i++) {
					var columns = get_columns(__DATA__, chartTypes[i].type, chartTypes[i].getter);
					draw(chartTypes[i].type, chartTypes[i].label, [x].concat(columns));
				}
			})();
  	</script>
	</body>
</html>  	
`
