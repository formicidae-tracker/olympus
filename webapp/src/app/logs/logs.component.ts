import { Component, OnInit } from '@angular/core';
import { ServiceLog, ServiceLogReport } from '@models/service-log';
import { OlympusService } from '@services/olympus';

@Component({
	selector: 'app-logs',
	templateUrl: './logs.component.html',
	styleUrls: ['./logs.component.css']
})
export class LogsComponent implements OnInit {

	public report: ServiceLogReport;

	public loopData: any[] = [
		{title: "Tag Trackers",section: "trackings"},
		{title: "Climates",section: "climates"},
	];
	constructor(private olympus: OlympusService) {
		this.report = new ServiceLogReport();
	}

	ngOnInit(): void {
		this.olympus.logs().subscribe(
			(report) => {
				this.report = report;
			});
	}

	displayState(l: ServiceLog): string {
		if ( l.on == true ) {
			return 'On';
		}
		if ( l.graceful == true ) {
			return 'Off';
		}
		return 'Exited prematurely';
	}

	actionState(l: ServiceLog,i: number): string {
		if ( l.on == false )
			if (l.graceful == false ) {
				return 'warning';
			} else {
				return 'info';
		}
		return '';
	}


}
