import { Component, Input, OnInit } from '@angular/core';
import { formatDate } from '@angular/common';
import { ClimateTimeSeries } from '@models/climate-time-series';
import { Chart,ChartDataset, ChartData, ChartOptions, ChartType } from 'chart.js';

@Component({
	selector: 'app-climate-chart',
	templateUrl: './climate-chart.component.html',
	styleUrls: ['./climate-chart.component.css']
})
export class ClimateChartComponent implements OnInit {
	@Input() units: string;

	public chartOptions: ChartOptions = {
		responsive: true,
		animation: null,
		plugins: {
			legend: {
				position: 'bottom',
			},
			tooltip: {
				callbacks: {
					label: function(context) {
						// let label: string = 'Time: '

						// if (context.parsed.x != null ) {
						// 	label += formatDate(new Date(new Date().getTime() + parseFloat(tooltipItem.label) * 3600000),'YYYY-MM-dd HH:mm:ss','en-US')+', '

						// label += data.datasets[tooltipItem.datasetIndex].label + ': ';

						// label += Math.round(parseFloat(tooltipItem.value) * 100) / 100;
						// return label;
						let label: string = 'Time: todo';
						return label
					},
				},
			},
		},
		scales: {
			xTime: {
				type: 'linear',
				display: true,
				title: {
					display: true,
					text: 'Time (h)',
				},
			},
			yTemperature: {
				type: 'linear',
				display: true,
				position: 'left',
				grid: {
					drawOnChartArea: true,
				},
				title:{
					display: true,
					text: 'Temperature (°C)',
				},
				suggestedMin: 0,
				suggestedMax: 25,

			},
			yHumidity: {
				type: 'linear',
				display: true,
				position: 'right',
				grid: {
					drawOnChartArea: false,
				},
				title:{
					display: true,
					text: 'Relative Humidity (%)',
				},
				suggestedMin: 0,
				suggestedMax: 90,
			},
		},
	};


	public chartData: ChartDataset[] = [];


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
			  color: string): ChartDataset {
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
