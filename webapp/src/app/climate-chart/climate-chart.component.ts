import { Component, Input, OnInit } from '@angular/core';
import { formatDate } from '@angular/common';
import { ClimateTimeSeries } from '@models/climate-time-series';
import { Chart, ChartDataset, ChartData, ChartOptions, ChartType, CartesianScaleOptions } from 'chart.js';

@Component({
	selector: 'app-climate-chart',
	templateUrl: './climate-chart.component.html',
	styleUrls: ['./climate-chart.component.css']
})


export class ClimateChartComponent implements OnInit {
	@Input() units: string;


	public chartOptions: ChartOptions = {
		responsive: true,
		animation: false,
		plugins: {
			legend: {
				position: 'bottom',
			},
			tooltip: {
				callbacks: {
					label: function(context) {
						if (context.dataset.data.length == 0 ) {
							return '';
						}
						let reference: Date = new Date(context.dataset.data[0]['reference']);
						let ratio: number = context.dataset.data[0]['ratio']
						let unit: string = context.dataset.data[0]['unit']
						let time: Date = new Date(reference.getTime() +  1000.0 * ratio * context.parsed.x)
						return time.toLocaleString() + ': ' + context.parsed.y + ' ' + unit;
					},
				}
			}
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
		(this.chartOptions.scales.xTime as CartesianScaleOptions).title.text = "Time (" + value.units + ")";
		if ( value.hasTemperature() == false ) {
			this.chartOptions.scales.yTemperature.display = false;
			this.chartOptions.scales.yHumidity.display = true;
			this.chartOptions.scales.yHumidity.grid.drawOnChartArea = true;
			(this.chartOptions.scales.yHumidity as CartesianScaleOptions).position = 'left';
		} else {
			this.chartOptions.scales.yTemperature.display = true;
			(this.chartOptions.scales.yHumidity as CartesianScaleOptions).position = 'right';
			this.chartOptions.scales.yHumidity.display = value.hasHumidity();
			this.chartOptions.scales.yHumidity.grid.drawOnChartArea = false;
		}
	}

	private colors = {
		humidity: '#1f77b4',
		temperature: '#ff7f0e',
		aux: [ '#2ca02c','#17becf','#ff6384']
	};

	private unitsToSeconds = new Map<string,number>([
		["m",60],
		["h",3600],
		["d",24 * 3600],
	]);

	buildData(series: ClimateTimeSeries,
			  values: any[],
			  name: string,
			  axis: string,
			  unit: string,
			  color: string): ChartDataset {
		if ( values.length > 0) {
			values[0].reference = series.reference;
			values[0].ratio = this.unitsToSeconds.get(series.units) ||  3600;
			values[0].unit = unit;
		}
		return {
			label: name,
			borderColor: color,
			fill: false,
			showLine: true,
			data: values,
			yAxisID: axis,
			pointRadius: 0,
			pointBackgroundColor: color,
			pointHoverBackgroundColor: color,
			backgroundColor: color,
			parsing: false,
		};
	}

	updateData(series: ClimateTimeSeries): void {
		this.chartData = [];
		if ( series.humidity.length > 0 ) {
			this.chartData.push(this.buildData(series,
											   series.humidity,
											   'Humidity',
											   'yHumidity',
											   '% R.H.',
											   this.colors.humidity));
		}
		if ( series.temperature.length > 0 ) {
			this.chartData.push(this.buildData(series,
											   series.temperature,
											   'Temperature',
											   'yTemperature',
											   '°C',
											   this.colors.temperature));
		}
		let numAux = Math.min(3,series.temperatureAux.length)
		for ( let i = 0; i < numAux; i++ ) {
			this.chartData.push(this.buildData(series,
											   series.temperatureAux[i],
											   'Temperature Aux ' + i,
											   'yTemperature',
											   '°C',
											   this.colors.aux[i]));
		}
	}

	constructor() { }

	ngOnInit(): void {
	}

}
