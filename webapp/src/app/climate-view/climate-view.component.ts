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

	public hostName: string;
	public zoneName: string;
	public units: string;

	private unitsFromWindowMap: Map<string,string>;

	@Input()
	set host(h: string) {
		this.hostName = h;
		this.updateChart();
	}

	@Input()
	set zone(z: string) {
		this.zoneName = z;
		this.updateChart();
	}

	constructor(private olympus: OlympusService) {
		this.hostName = '';
		this.zoneName = '';
		this.window = '1d';
		this.units = 'h';
		this.unitsFromWindowMap = new Map<string,string>([
			["10m","m"],
			["1h","m"],
			["1d","h"],
			["1w","h"]
		]);
		this.climateTimeSeries = new ClimateTimeSeries();
	}

	ngOnInit() {
		this.update = timer(0,10000).subscribe(() => {
			this.updateChart();
		});
	}

	ngOnDestroy() {
		if ( this.update != null ) {
			this.update.unsubscribe();
		}
	}


	updateChart() {
		if ( this.hostName.length == 0
			 || this.zoneName.length == 0 ) {
			this.climateTimeSeries = new ClimateTimeSeries();
			return;
		}

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

	unitsFromWindow(window: string): string {
		if ( this.unitsFromWindowMap.has(window) ) {
			return this.unitsFromWindowMap.get(window);
		}
		return 'h';
	}


	public selectTimeWindow(window: string) {
		this.window = window;
		this.units = this.unitsFromWindow(window);
		this.updateChart();
	}



}
