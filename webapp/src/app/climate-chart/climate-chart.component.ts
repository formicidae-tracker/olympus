import { Component, Input, OnInit } from '@angular/core';
import { ClimateTimeSeries } from '@models/climate-time-series';
import { ChartDataSets, ChartOptions, ChartType } from 'chart.js';

@Component({
	selector: 'app-climate-chart',
	templateUrl: './climate-chart.component.html',
	styleUrls: ['./climate-chart.component.css']
})
export class ClimateChartComponent implements OnInit {

	public chartOptions: ChartOptions = {
		responsive: true,
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
					position: 'right',
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
					position: 'left',
					id: 'y-humidity',
					gridLines: {
						drawOnChartArea: false,
					},
					scaleLabel:{
						display: true,
						labelString: 'Relative Humidity (%)',
					},
					ticks: {
						suggestedMin: 20,
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
		if ( value.temperature.length == 0
			&& value.temperatureAux.length == 0 ) {
			this.chartOptions.scales.yAxes[0].display = false;
			this.chartOptions.scales.yAxes[1].display = true;
			this.chartOptions.scales.yAxes[1].gridLines.drawOnChartArea = true;
			this.chartOptions.scales.yAxes[1].position = 'left';
		} else {
			this.chartOptions.scales.yAxes[1].gridLines.drawOnChartArea = false;
			this.chartOptions.scales.yAxes[1].position = 'right';
			this.chartOptions.scales.yAxes[0].display = true;
			this.chartOptions.scales.yAxes[1].display = value.humidity.length > 0;
		}
	}

	private colors = {
		humidity: '#1f77b4',
		temperature: '#1f77b4',
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
			yAxisID: axis
		};
		let last = values[values.length-1].X;
		let timeDiv = 60;
		if ( last > 3600 ) {
			timeDiv = 3600;
			this.chartOptions.scales.xAxes[0].scaleLabel.labelString = 'Time (h)';
		} else {
			this.chartOptions.scales.xAxes[0].scaleLabel.labelString = 'Time (m)';
		}

		for ( let p of values ) {
			res.data.push({
				x: Math.round((last - p.X)/timeDiv*100)/100,
				y: Math.round(100*p.Y)/100,
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
			this.chartData.push(this.buildData(value.humidity,
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
