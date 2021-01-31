export enum AlarmLevel {
	Warning = 1,
	Critical = 2,
}

export class Alarm {

	constructor(public Reason :string,
				public On: boolean,
				public LastChange: Date,
				public Level: AlarmLevel,
				public Triggers: number
			   ) {}

	action() {
		if (this.On == false) {
			return 'info';
		}
		if (this.Level == AlarmLevel.Warning) {
			return 'warning';
		}
		return 'danger';
	}

	static adapt(item: any): Alarm {
		return new Alarm(
			item.Reason,
			item.On,
			item.LastChange==null?null:new Date(item.LastChange),
			item.Level,
			item.Triggers
		);
	}
}

export function CompareAlarm(a :Alarm, b :Alarm){
	if (a.On ==  b.On ) {
		if (a.Level < b.Level) {
			return 1;
		} else if (a.Level > b.Level) {
			return -1;
		}
		if (a.LastChange == b.LastChange ) {
			return a.Reason.localeCompare(b.Reason);
		}
		if (a.LastChange < b.LastChange ) {
			return 1;
		} else {
			return -1;
		}
	}
	if (a.On == true) {
		return -1;
	}
	return 1;
}