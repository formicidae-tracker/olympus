import { Component, Input, OnInit } from '@angular/core';
import { formatDate } from '@angular/common';
import { ClimateTimeSeries } from '@models/climate-time-series';
import { Chart,ChartDataSets, ChartData, ChartTooltipItem, ChartOptions, ChartType } from 'chart.js';

@Component({
	selector: 'app-climate-chart',
	templateUrl: './climate-chart.component.html',
	styleUrls: ['./climate-chart.component.css']
})
export class ClimateChartComponent implements OnInit {

	public chartOptions: ChartOptions = {
		responsive: true,
		animation: null,
		legend: {
			position: 'bottom',
		},
		tooltips: {
            callbacks: {
                label: function(tooltipItem: ChartTooltipItem,data: ChartData) {

                    let label: string = 'Time: ' + formatDate(new Date(new Date().getTime() + parseFloat(tooltipItem.label) * 3600000),'YYYY-MM-dd HH:mm:ss','en-US')+', '

					label += data.datasets[tooltipItem.datasetIndex].label + ': ';

                    label += Math.round(parseFloat(tooltipItem.value) * 100) / 100;
                    return label;
                },
            },
		},
		scales: {
			xAxes: [
				{
					scaleLabel: {
						display: true,
						labelString: 'Time(h)',
					},
					display: true,
				},
			],
			yAxes: [
				{
					type: 'linear',
					display: true,
					position: 'left',
					id: 'y-temperature',
					gridLines: {
						drawOnChartArea: true,
					},
					scaleLabel:{
						display: true,
						labelString: 'Temperature (°C)',
					},
					ticks: {
						suggestedMin: 0,
						suggestedMax: 25,
					},
				},
				{
					type: 'linear',
					display: true,
					position: 'right',
					id: 'y-humidity',
					gridLines: {
						drawOnChartArea: false,
					},
					scaleLabel:{
						display: true,
						labelString: 'Relative Humidity (%)',
					},
					ticks: {
						suggestedMin: 0,
						suggestedMax: 90,
					},
				},
			],
		}
	};


	public chartData: ChartDataSets[] = [];


	public chartType: ChartType = 'scatter';

	@Input()
	set climateTimeSeries(value: ClimateTimeSeries) {
		this.updateChartOptions(value);
		this.updateData(value);
	}

	updateChartOptions(value: ClimateTimeSeries): void {
		if (value.hasData() == false ) {
			return;
		}

		if ( value.hasTemperature() == false ) {
			this.chartOptions.scales.yAxes[0].display = false;
			this.chartOptions.scales.yAxes[1].display = true;
			this.chartOptions.scales.yAxes[1].gridLines.drawOnChartArea = true;
			this.chartOptions.scales.yAxes[1].position = 'left';
		} else {
			this.chartOptions.scales.yAxes[0].display = true;
			this.chartOptions.scales.yAxes[1].position = 'right';
			this.chartOptions.scales.yAxes[1].display = value.hasHumidity();
			this.chartOptions.scales.yAxes[1].gridLines.drawOnChartArea = false;
		}
	}

	private colors = {
		humidity: '#1f77b4',
		temperature: '#ff7f0e',
		aux: [ '#2ca02c','#17becf','#ff6384']
	};

	buildData(values: any[],
					 name: string,
					 axis: string,
					 color: string): ChartDataSets {
		let res = {
			label: name,
			borderColor: color,
			fill: false,
			showLine: true,
			lineTension: 0,
			data: [],
			yAxisID: axis,
			pointRadius: 0,
			pointBackgroundColor: color,
			pointHoverBackgroundColor: color,
		};
		let last = values[values.length-1].X;
		let timeDiv = 3600;

		for ( let p of values ) {
			res.data.push({
				x: (p.X - last)/timeDiv,
				y: p.Y,
			});
		}
		return res;
	}

	updateData(value: ClimateTimeSeries): void {
		this.chartData = [];
		if ( value.humidity.length > 0 ) {
			this.chartData.push(this.buildData(value.humidity,
											   'Humidity',
											   'y-humidity',
											   this.colors.humidity));
		}
		if ( value.temperature.length > 0 ) {
			this.chartData.push(this.buildData(value.temperature,
											   'Temperature',
											   'y-temperature',
											   this.colors.temperature));
		}
		let numAux = Math.min(3,value.temperatureAux.length)
		for ( let i = 0; i < numAux; i++ ) {
			this.chartData.push(this.buildData(value.temperatureAux[i],
											   'Temperature Aux ' + i,
											   'y-temperature',
											   this.colors.aux[i]));
		}
	}

	constructor() { }

	ngOnInit(): void {
	}

}
