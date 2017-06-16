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
		<script src="https://cdnjs.cloudflare.com/ajax/libs/d3/3.5.17/d3.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/crossfilter/1.3.12/crossfilter.min.js"></script>
		<script src="http://cdnjs.cloudflare.com/ajax/libs/dc/2.1.6/dc.min.js"></script>
		<link rel="stylesheet" href="http://cdnjs.cloudflare.com/ajax/libs/dc/2.1.6/dc.css"/>
	</head>
  	<body>
	  	<div id="memory-chart"></div>
		<div id="cpu-chart"></div>
		<div id="keys-chart"></div>
		<div id="inserts-chart"></div>
		<div id="reads-chart"></div>

	  	<script>
		  /*
		  	var __DATA__ = [
				{
					name: 'f1',
					address: '172.0.0.1',
					stats: [
						{
							when: '2017-06-16T14:01:33',
							memory: 2,
						},
						{
							when: '2017-06-16T14:01:33',
							memory: 2,
						},
					]
				},
				{
					name: 'f2',
					address: '172.0.0.2',
					stats: [
						{
							when: '2017-06-16T14:01:33',
							memory: 2,
						},
						{
							when: '2017-06-16T14:01:33',
							memory: 2,
						},
					]
				}
			];
			*/

			var __DATA__ = [
				{
					date: '2017-06-16T14:01:33',
					stats: {
						'f1': {
							memory: 2
						},
						'f2': {
							memory: 2
						}
					}
				},
				{
					date: '2017-06-16T14:01:33',
					stats: {
						'f1': {
							memory: 2
						},
						'f2': {
							memory: 2
						},
						'f3': {
							memory: 2
						}
					}
				}
			];

		  	(function() {
				var memoryChart = dc.lineChart('#memory-chart');
				var cpuChart = dc.lineChart('#cpu-chart');
				var keysChart = dc.lineChart('#keys-chart');
				var insertsChart = dc.lineChart('#insets-chart');
				var readsChart = dc.lineChart('#reads-chart');

				var ndx = crossfilter(__DATA__);

				var dimension = ndx.dimension(function(d) {
        			return new Date(d.date);
    			});
				var nodes = [];

				

				var nodeGroup = dimension.group();
				var memoryGroup = dimension.group().reduce(function(p, v) {
                  	for (var key in v.stats) {
  						if (v.stats.hasOwnProperty(key)) {
    						p[key] = v.stats[key].memory;
  						}
					}
				  	return p;
              	}, function(p, v) {
                  	return p;
              	}, function() {
                	return {};
              	});

				function for_node(node_name) {
					return function(d) {
						return d.value[node_name];
					}
				}

				memoryChart
					.renderArea(true)
					.width(990)
					.height(200)
					.transitionDuration(1000)
					.margins({
						top: 30,
						right: 50,
						bottom: 25,
						left: 40
					})
					.x(d3.time.scale().domain([d3.time.hour.offset(new Date(), -1), new Date()]))
					.dimension(dimension)
					.group(memoryGroup, __DATA__[0].name, for_node(__DATA__[0].name));
				for(var i = 1; i < __DATA__.length; ++i) {
              		memoryChart.stack(memoryGroup, __DATA__[i].name, for_node(__DATA__[i].name));
				}
          		memoryChart.render();
			})();
	  	</script>
	</body>
</html>  	
`
