import { formatDate,formatNumber,formatPercent } from '@angular/common';
import { Component, OnInit, Input } from '@angular/core';
import { ZoneSummaryReport } from '@models/zone-summary-report';


@Component({
	selector: 'app-zone-preview',
	templateUrl: './zone-preview.component.html',
	styleUrls: ['./zone-preview.component.css']
})

export class ZonePreviewComponent implements OnInit {
    @Input() summary: ZoneSummaryReport;

    constructor() {
		this.summary = new ZoneSummaryReport('','');
	}

    ngOnInit() {
    }


	detailRouterLink(): string[] {
		return ['host',this.summary.host,'zone',this.summary.zoneName];
	}

	zoneDisplayName(): string {
		return this.summary.host + '.' + this.summary.zoneName;
	}

	currentStateDescription(): string {
		if ( this.summary.climate == null
		    || this.summary.climate.current == null) {
			return 'No climate control';
		}
		return 'Current state: \'' + this.summary.climate.current.name + '\'';
	}

	nextStateDescription(): string {
		if ( this.summary.climate == null ) {
			return '';
		}
		if ( this.summary.climate.next == null ) {
			return 'No next state';
		}
		return 'Next state: \'' + this.summary.climate.next.name + '\'';
	}

	hasNextTime(): boolean {
		return this.summary.climate != null && this.summary.climate.nextTime != null;
	}
	nextTime(): Date {
		return this.summary.climate.nextTime;
	}

	hasTemperature(): boolean {
		return this.summary.climate != null && isNaN(this.summary.climate.temperature) == false;
	}


	temperature(): string {
		if ( this.summary.climate == null ) {
			return 'N.A.';
		}
		return Math.round(100*this.summary.climate.temperature)/100+' °C';
	}

	temperatureStatus(): string {
		if ( this.hasTemperature() == false ) {
			return '';
		}
		return this.summary.climate.temperatureBounds.status(this.summary.climate.temperature);
	}

	hasHumidity(): boolean {
		return this.summary.climate != null && isNaN(this.summary.climate.humidity) == false;
	}

	humidity(): string {
		if ( this.summary.climate == null ) {
			return 'N.A.';
		}
		return Math.round(100*this.summary.climate.humidity)/100+' %';
	}

	humidityStatus(): string {
		if ( this.hasHumidity() == false ) {
			return '';
		}
		return this.summary.climate.humidityBounds.status(this.summary.climate.humidity);
	}

	alarmStatus(): string {
		if ( this.summary.climate == null ) {
			return '';
		}
		if ( this.summary.climate.activeEmergencies > 0 ) {
			return 'danger';
		}
		if ( this.summary.climate.activeWarnings > 0 ) {
			return 'warning';
		}
		return 'info';
	}

	activeEmergencies(): number {
		if ( this.summary.climate == null ) {
			return 0;
		}
		return this.summary.climate.activeEmergencies;
	}

	activeWarnings(): number {
		if ( this.summary.climate == null ) {
			return 0;
		}
		return this.summary.climate.activeWarnings;
	}

}
