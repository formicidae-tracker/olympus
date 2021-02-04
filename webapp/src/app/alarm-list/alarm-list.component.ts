import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { AlarmReport } from '@models/alarm';
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
	collapsed: Map<string,boolean>;


	constructor() {
		this.now = new Date();
		this.alarms = [];
	}



	ngOnInit(): void {
		this.updateTime = timer(0,1000).subscribe(() => {
			this.now = new Date();
		})
		this.collapsed = new Map<string,boolean>();
	}

	ngOnDestroy(): void {
		if ( this.updateTime != null ) {
			this.updateTime.unsubscribe();
		}
	}

	toggleCollapse(r: string): void {
		if ( this.collapsed.has(r) == false ) {
			this.collapsed.set(r,true);
		}
		this.collapsed.set(r,!this.collapsed.get(r));
	}

	isCollapsed(r: string): boolean {
		if ( this.collapsed.has(r) == false ) {
			return true;
		}

		return this.collapsed.get(r);
	}

	disabled(): boolean {
		for ( let r of this.alarms ) {
			if ( r.on() == true ) {
				return false;
			}
		}
		return true;
	}

}
