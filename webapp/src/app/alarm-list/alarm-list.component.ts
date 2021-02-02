import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { environment } from '@environments/environment';
import { AlarmEvent, AlarmReport } from '@models/alarm';
import { Subscription,timer } from 'rxjs';

@Component({
	selector: 'app-alarm-list',
	templateUrl: './alarm-list.component.html',
	styleUrls: ['./alarm-list.component.css']
})
export class AlarmListComponent implements OnInit,OnDestroy {
	@Input() alarms: AlarmReport[];
	now: Date;
	updateTime: Subscription;

	constructor() {
		this.now = new Date();
		if ( environment.production == true ) {
			this.alarms = [];
		} else {
			let now = (new Date()).getTime();
			this.alarms = [
				new AlarmReport('Water level is critical',
								2,
								[
									new AlarmEvent(true,new Date(now - 30*60000)),
								]),
				new AlarmReport('Humidity is unreachable',
								1,
								[
									new AlarmEvent(true,new Date(now - 60*60000)),
									new AlarmEvent(false,new Date(now - 50*60000)),
									new AlarmEvent(true,new Date(now - 40*60000)),
								]),
				new AlarmReport('Temperature is out of bound',
								2,
								[
									new AlarmEvent(true,new Date(now - 65*60000)),
									new AlarmEvent(false,new Date(now - 45*60000)),
								]),
			]
		}
	}

	ngOnInit(): void {
		this.updateTime = timer(0,1000).subscribe(() => {
			this.now = new Date();
		})
	}

	ngOnDestroy(): void {
		if ( this.updateTime != null ) {
			this.updateTime.unsubscribe();
		}
	}
}
