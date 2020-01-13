function ppastats_chart(ddts) {
	var data_chart = [];
	var max_date = null;
	var min_date = null;

	$.each(ddts, function(i, item) {
		var tm = item["time"];
		var d = new Date(tm[0], tm[1]-1, tm[2]);
		var entry = [d, item["value"]];
		data_chart.push(entry);

		if (max_date == null || max_date < d) {
			max_date = d;
		}

		if (min_date == null || min_date > d) {
			min_date = d;
		}
	});

	var plot1 = $.jqplot ('chart', [data_chart], {
		title: 'Daily Download Count',
		axes: {
			xaxis: {
				renderer:$.jqplot.DateAxisRenderer,
				tickOptions:{formatString:'%Y/%m/%d'},
				min: min_date,
				max: max_date
			},
			yaxis: {
				min: 0
			}
		},
		series: [{lineWidth:1,showMarker:false}]
	});
}

function ppastats_distros(distros) {
	var data_chart = [];
	var max_date = null;
	var min_date = null;
	var series_opt = [];

	$.each(distros, function(i, distro) {
		var arr = [];

		$.each(distro["ddts"], function(j, item) {
			var tm = item["time"];
			var d = new Date(tm[0], tm[1]-1, tm[2]);
			var entry = [d, item["value"]];
			arr.push(entry);

			if (max_date == null || max_date < d) {
				max_date = d;
			}

			if (min_date == null || min_date > d) {
				min_date = d;
			}
		});
		data_chart[i] = arr;
		series_opt[i] = {label: distro["name"]};
	});

	var plot1 = $.jqplot ('chart_distro', data_chart, {
		title: 'Daily Download Count by Ubuntu Distribution',
		axes: {
			xaxis: {
				renderer:$.jqplot.DateAxisRenderer,
				tickOptions:{formatString:'%Y/%m/%d'},
				min: min_date,
				max: max_date
			},
			yaxis: {
				min: 0
			}
		},
		seriesDefaults: {
			lineWidth:1,
			showMarker:false
		},
		legend: {
			show: true
		},
		series: series_opt
	});
}

function ppastats_pkg(json_url) {
	$(document).ready(function() {
		$.getJSON(json_url, function(data) {
			var downloads = 0;

			$("#ppa_owner").html(data["ppa_owner"]);
			$("#ppa_name").html(data["ppa_name"]);
			$("#pkg_name").html(data["name"]);

			$("#versions").append("<ul>");
			$.each(data["versions"], function(i, v) {
				var v_url = data["name"]+"_"+v+".html";

				$("#versions ul").append("<li><a href='"+v_url+"'>"+v+"</a></li>");
			});

			$("#distros").append("<ul>");
			$.each(data["distros"], function(i, d) {
				$("#distros ul").append("<li>"+d["name"]+": "+d["count"]+"</li>");
				downloads += d["count"];
			});

			$("#pkg_downloads").html("" + downloads);

			ppastats_chart(data["ddts"]);
			ppastats_distros(data["distros"]);
		});
	});
}

function ppastats_ver() {
	$(document).ready(function() {
		var pkg_url = data["pkg_name"]+".html";

		$("#ppa_owner").html(data["ppa_owner"]);
		$("#ppa_name").html(data["ppa_name"]);
		$("#pkg_name").html("<a href='"+pkg_url+"'>"+data["pkg_name"]+"</a>");
		$("#version").append(" "+data["name"]);

		$("#distros").append("<ul>");
		$.each(data["distros"], function(i, distro) {
			$.each(distro["archs"], function(i, arch) {
				$("#distros ul").append("<li>"+distro["name"]+"_"+arch["name"]+": "+arch["count"]+"</li");
			});
		});

		ppastats_chart(data["ddts"]);

		$("#date_created").append(data["date_created"]);
	});
}

function ppastats_ppa() {
    $(document).ready(function() {
	var max_date = null;
	var min_date = null;
	var pkg_url;

	$.getJSON("index.json", function(data) {
	    pkg_url = data["pkg_name"]+".html";

	    $("#ppa_name").html(data["ppa_owner"]+"/"+data["ppa_name"]);

	    $.each(data["packages"], function(i, item) {
		var url = item["name"]+".html";
		$("#pkgs").append("<li><a href='"+url+"'>"+item["name"]+"</a>: "+item["count"]+"</li>");
	    });

	    ppastats_chart(data["ddts"]);
	});
    });
}
