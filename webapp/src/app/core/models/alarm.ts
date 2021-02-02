export enum AlarmLevel {
	Warning = 1,
	Critical = 2,
}

export class AlarmEvent {
	constructor(public on: boolean = false,
				public time: Date = new Date()) {
	}
	static adapt(item: any): AlarmEvent {
		return new AlarmEvent(item.On,new Date(item.Time));
	}

}

export class AlarmReport {
	public count: number;
	constructor(public reason: string = '',
				public level: AlarmLevel = AlarmLevel.Warning,
				public events: AlarmEvent[] = []) {
		this.count = 0;
		for ( let e of events ) {
			if ( e.on == true ) {
				this.count++;
			}
		}
	}

	static adapt(item: any): AlarmReport{
		let events: AlarmEvent[] = [];
		for ( let e of item.Events ) {
			events.push(AlarmEvent.adapt(e));
		}
		return new AlarmReport(item.Reason,item.Level,events);
	}

	lastEvent(): AlarmEvent {
		if ( this.events.length == 0 ) {
			return undefined;
		}

		return this.events[this.events.length-1];
	}

	on(): boolean {
		return this.lastEvent().on || false;
	}

	lastTime(): Date {
		return this.lastEvent().time || new Date();
	}

	action() {
		if (this.on() == false) {
			return 'info';
		}
		if (this.level == AlarmLevel.Warning) {
			return 'warning';
		}
		return 'danger';
	}

	since(now: Date): string {
		let ellapsed = now.getTime() - this.lastTime().getTime();
		if ( ellapsed <= 1000 ) {
			return 'now';
		}
		ellapsed = Math.round(ellapsed/1000);
		let seconds = ellapsed%60;
		if (ellapsed < 60 ) {
			return seconds+'s';
		}
		let minutes = Math.floor(ellapsed/60)
		if ( minutes < 60 ) {
			return minutes+'m';
		}
		let hours = Math.floor(minutes/60)
		minutes = minutes % 60;
		return hours+'h'+minutes+'m';
	}


	static compare(a: AlarmReport, b: AlarmReport): number {
		let aOn = a.on();
		let bOn = b.on();
		if ( aOn == bOn ) {
			if ( a.level == b.level ) {
				let atime = a.lastTime();
				let btime = b.lastTime();
				if ( atime == btime ) {
					return a.reason.localeCompare(b.reason);
				}
				return (atime < btime) ? -1 : 1;
			}
			return (a.level < b.level) ? -1 : 1;
		}
		return aOn ? -1 : 1;
	}

}
