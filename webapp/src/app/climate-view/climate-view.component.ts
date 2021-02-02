import { Component, OnInit, OnDestroy, Input} from '@angular/core';
import { ClimateTimeSeries } from '@models/climate-time-series';
import { OlympusService } from '@services/olympus';
import { Subscription, timer } from 'rxjs';


@Component({
	selector: 'app-climate-view',
	templateUrl: './climate-view.component.html',
	styleUrls: ['./climate-view.component.css']
})

export class ClimateViewComponent implements OnInit,OnDestroy {

	public window: string;
	public climateTimeSeries: ClimateTimeSeries;
	update: Subscription;


	@Input() hostName: string;
	@Input() zoneName: string;

	constructor(private olympus: OlympusService) {
		this.hostName = '';
		this.zoneName = '';
		this.window = '1d';
		this.climateTimeSeries = new ClimateTimeSeries();
	}

	ngOnInit() {
		this.update = timer(0,10000).subscribe(  () => {
			this.updateChart();
		});
	}

	ngOnDestroy() {
		if ( this.update != null ) {
			this.update.unsubscribe();
		}
	}


	updateChart() {
		this.olympus.climateTimeSeries(this.hostName,this.zoneName,this.window)
			.subscribe((timeSeries) => {
				this.climateTimeSeries = timeSeries;
			},() => {
				this.climateTimeSeries = new ClimateTimeSeries();
			});

	}

	isSelected(window: string) {
		if ( window == this.window ) {
			return ' active'
		}
		return ''
	}

	public selectTimeWindow(window: string) {
		this.window = window;
		this.updateChart();
	}



}
