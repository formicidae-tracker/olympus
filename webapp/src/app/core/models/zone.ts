import { State } from '@models/state';
import { Alarm,AlarmLevel,CompareAlarm} from '@models/alarm';
import { Bounds } from './bounds';


export class Zone {
	constructor(public Host: string,
				public Name: string,
				public Temperature: number,
				public TemperatureBounds: Bounds,
				public Humidity: number,
				public HumidityBounds: Bounds,
				public Alarms: Alarm[],
				public Current: State,
				public CurrentEnd: State,
				public Next: State,
				public NextEnd: State,
				public NextTime: Date) {}

    temperatureStatus() {
		if (this.Temperature < this.TemperatureBounds.Min ) {
			return 'warning';
		}
		if (this.Temperature > this.TemperatureBounds.Max ) {
			return 'danger';
		}
		return 'success';
    }

	FullName() {
		return this.Host + '.' + this.Name;
	}

    humidityStatus() {
		if (this.Humidity < this.HumidityBounds.Min ) {
			return 'danger';
		}
		if (this.Humidity > this.HumidityBounds.Max ) {
			return 'warning';
		}
		return 'success';
    }

    alarmStatus() {
		if (this.Alarms.length == 0 ) {
			return 'danger';
		}
		let aRes = this.Alarms[0];
		for ( let a of this.Alarms ) {
			if (a.On == false ) {
				continue;
			}
			if ( aRes.On == false || a.Level > aRes.Level ) {
				aRes = a;
			}

		}
		return aRes.action();
    }

	sortedAlarm() {
		let res = Object.assign([],this.Alarms);
		return res.sort(CompareAlarm);
	}

	numberOfActiveAlarms() {
		let res = 0;
		for ( let a of this.Alarms ) {
			if (a.On == true ) {
				res += 1;
			}
		}
		return res;
	}

	static adapt(item: any): Zone {
		let alarms: Alarm[] = [];
		if (item.Alarms != null) {
			for ( let a of item.Alarms ) {
				alarms.push(Alarm.adapt(a));
			}
		}
		let current: State = null;
		let currentEnd: State = null;
		let next: State = null;
		let nextEnd: State = null;
		let nextTime: Date = null;
		if (item.Current != null) {
			current = State.adapt(item.Current);
		}
		if (item.CurrentEnd != null ) {
			currentEnd = State.adapt(item.CurrentEnd);
		}

		if ( item.Next != null && item.NextTime != null ) {
			next = State.adapt(item.Next);
			nextTime = new Date(item.NextTime);
		}

		if ( item.NextEnd != null ) {
			nextEnd = State.adapt(item.NextEnd);
		}

		return new Zone(
			item.Host,
			item.Name,
			item.Temperature,
			Bounds.adapt(item.TemperatureBounds),
			item.Humidity,
			Bounds.adapt(item.HumidityBounds),
			alarms,
			current,
			currentEnd,
			next,
			nextEnd,
			nextTime
		);
	}
}
