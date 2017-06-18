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
				background: #636387;
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
				var output = [],
					columns = {};

				for (var i = 0; i < data.length; i++) {
					if (columns['x'] == undefined) {
						columns['x'] = [];
					} 
					columns['x'].push(new Date(data[i].date));
					for (var j = 0; j < data[i].stats.length; j++) {
						var node = data[i].stats[j].node;
						var val = data[i].stats[j][type];
						if (columns[node] == undefined) {
							columns[node] = [];
						}
						columns[node].push(getter? getter(val) : val); 
					}
				}

				var keys = Object.keys(columns);
				for (var i = 0; i < keys.length; i++) {
					var key = keys[i];
					output.push([key].concat(columns[key]));
				}
				return output;
			}

		  	function draw(type, label, getter) {
				var columns = get_columns(__DATA__, type, getter);  
				var dataTypes = {};
				var dataGroups = [];

				for (var i = 1; i < columns.length; i++) {
					dataGroups.push(columns[i][0]);
				}

				for (var i = 0; i < dataGroups.length; i++) {
					dataTypes[dataGroups[i]] = 'area-spline';
				}

				c3.generate({
					bindto: '#' + type + '-chart',
					data: {
						x: 'x',
						columns: columns,
						//types: dataTypes,
						//groupd: dataGroups
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
				var keys = [],
					sum = 0;
					nodes = __DATA__[__DATA__.length - 1].stats;

				for (var i = 0; i < nodes.length; i++) {
					var n = nodes[i].node,
						v = nodes[i].keys;
					sum = sum + v;
					keys.push([n, v]);
				}
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
				for (var i = 0; i < chartTypes.length; i++) {
					var type = chartTypes[i].type;
					var label = chartTypes[i].label;
					var getter = chartTypes[i].getter;

					draw(type, label, getter);
				}
			})();
  	</script>
	</body>
</html>  	
`
