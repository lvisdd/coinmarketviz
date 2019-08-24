var margin = {top: 30, right: 20, bottom: 30, left: 50},
    width = 600 - margin.left - margin.right,
    height = 400 - margin.top - margin.bottom;

var svg = d3.select("#ticker")
    .append("svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom);
    // .append("g")
    //     .attr("transform", 
    //           "translate(" + margin.left + "," + margin.top + ")");

var zoomLayer  = svg.append("g");
var zoomed = function() {
  zoomLayer.attr("transform", d3.event.transform);
}

svg.call(d3.zoom()
  .scaleExtent([1/2, 120])
  .on("zoom", zoomed)); 

// var url = "https://api.coinmarketcap.com/v1/ticker/"
// d3.json("test.json", function(error, json) {
var url = "/ticker/"
d3.json(url, function(error, json) {
  var width = document.querySelector("svg").clientWidth;
  var height = document.querySelector("svg").clientHeight;
  var div = d3.select("#table");
  var data =
  {
    "id": "coinmarketcap",
    "name": "CoinMarketCap",
    "children": json
  }

  root = d3.hierarchy(data);
  root.sum(function(d) { return d["24h_volume_usd"]; });

  var pack = d3.pack()
    .size([width, height])
    .padding(0);

  pack(root);

  var selection = d3.select("#ranking").selectAll(".table")
                    .append("tbody")
                    .data(root.descendants())
                    .enter()
                    .append("tr")
                    .html(function (d) {
                      if(d.children) return; 
                      return ranking(d.data);
                    });

  // var node = d3.select("svg").selectAll(".node")
  var node = zoomLayer.selectAll(".node")
    .data(root.descendants())
    // .data(root.leaves())
    .enter()
    .append("g")
    .attr("transform", function(d) { return "translate(" + d.x + "," + (d.y) + ")"; });
  
  var color = ["orange", "Khaki", "Ivory"];
  // var color = d3.scaleOrdinal(d3.schemeCategory20c);
  var text_color = [];
  text_color.push(d3.rgb(246, 53, 56)); // #F63538
  text_color.push(d3.rgb(174, 66, 72)); // #AE4248
  text_color.push(d3.rgb(111, 47, 51)); // #6F2F33
  text_color.push(d3.rgb(47, 50, 61));  // #2F323D
  text_color.push(d3.rgb(59, 90, 80));  // #3B5A50
  text_color.push(d3.rgb(49, 144, 78)); // #31904E
  text_color.push(d3.rgb(48, 190, 86)); // #30BE56

  // var colorScale = d3.scaleLinear()
  var colorScale = d3.scaleQuantize()
                     .domain([-15, 15])
                     .range(["#F63538", "#AE4248", "#6F2F33", "#2F323D", "#3B5A50", "#31904E", "#30BE56"]);
                     // .range([0, 6]);
                     // .range(['red','yellow','green'])
                     // .clamp(true);

  node.append("circle")
    .attr("r", function(d) { return d.r; })
    .attr("stroke", "black")
    .attr("fill", function(d) { 
      if(d.children) return "orange";
      return colorScale(d.data.percent_change_24h);
    })
    .attr("id", function(d) { return color[d.data.id]; })
    .on("mouseover", function(d) {
      if(d.children) return; 
      div.transition()
        .duration(200)
        .style("opacity", .9);
      div.html(tooltip(d.data))
        .style("left", (d3.event.pageX) + "px")
        .style("top", (d3.event.pageY - 28) + "px");
    })
    .on("mouseout", function(d) {
      if(d.children) return; 
      div.transition()
        .duration(200)
        .style("opacity", 0);
    })
    ;

  node.append("text")
    .style("text-anchor", function(d) { return d.children ? "end" : "middle"; })
    .text(function(d) { return d.children ? "" : d.data.name; })
    .attr("font-size", "75%")
    .attr("fill", "white")
    .each(function(d) {
        var bbox = this.getBBox();
        d.width = bbox.width;
        d.height = bbox.height;
        d.x = bbox.x;
        d.y = bbox.y;
    })
    .attr("font-size", function(d) { return d.height > (d.r / 4) ? (d.r / 4) : d.height; })
    ;
});

function tooltip(data) {
  var div = '<table class="table table-striped table-bordered table-sm w-auto">' + 
            '<thead class="mdb-color lighten-4">' +
            '<tr><th>Key</th><th class="th-lg">Value</th></tr>' + 
            '<tbody>';
  for (let k in data) {
    div = div + '<tr><th align="left">' + k + '</th><td>' + data[k] + '</td></tr>';
  }
  div = div + "</tbody></table>";
  return div;
};

function ranking(data) {
  var div = '<th align="left">' + data["rank"] + '</th><td>' + data["name"] + '</td>';
  return div;
};
