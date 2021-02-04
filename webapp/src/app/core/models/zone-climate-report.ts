import { State } from './state';
import { Bounds } from './bounds';

export class ZoneClimateReport {
	constructor(public temperature: number = NaN,
				public humidity: number = NaN,
				public temperatureBounds: Bounds = new Bounds(NaN,NaN),
				public humidityBounds: Bounds = new Bounds(NaN,NaN),
				public activeWarnings: number = 0,
				public activeEmergencies: number = 0,
				public numAux: number = 0,
				public current: State = null,
				public currentEnd: State = null,
				public next: State = null,
				public nextEnd: State = null,
				public nextTime: Date = null ){
		if ( this.temperature <= -1000.0 ) {
			this.temperature = NaN;
		}
		if ( this.humidity <= -1000.0 ) {
			this.humidity = NaN
		}
	}

	static adapt(item: any): ZoneClimateReport {
		if ( item == null ) {
			return null;
		}
		let nextTime: Date = null;
		if ( item.NextTime != null ) {
			nextTime = new Date(item.NextTime);
		}
		return new ZoneClimateReport(item.Temperature,
									 item.Humidity,
									 Bounds.adapt(item.TemperatureBounds),
									 Bounds.adapt(item.HumidityBounds),
									 item.ActiveWarnings,
									 item.ActiveEmergencies,
									 item.NumAux,
									 State.adapt(item.Current),
									 State.adapt(item.CurrentEnd),
									 State.adapt(item.Next),
									 State.adapt(item.NextEnd),
									 nextTime);

	}
}
