import { Component, OnInit, Input } from '@angular/core';

import { AlarmReport } from '@models/alarm';

@Component({
	selector: 'app-alarm',
	templateUrl: './alarm.component.html',
	styleUrls: ['./alarm.component.css']
})

export class AlarmComponent implements OnInit {

	@Input() alarm: AlarmReport;
	@Input() now: Date;
	public isCollapsed: boolean;
	constructor() {
		this.alarm = new AlarmReport();
		this.now = new Date();
		this.isCollapsed = true;
	}


	ngOnInit() {
	}

}
