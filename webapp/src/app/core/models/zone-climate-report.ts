
import { ZoneClimateStatus } from './zone-climate-status';
import { State } from './state';

export class ZoneClimateReport {
	constructor(public ClimateStatus: ZoneClimateStatus,
				public NumAux: number,
				public Current: State,
				public CurrentEnd: State,
				public Next: State,
				public NextEnd: State,
				public NextTime: Date){}

	static adapt(item: any): ZoneClimateReport {
		let nextTime: Date = null;
		if ( item.NextTime != null ) {
			nextTime = new Date(item.NextTime);
		}
		return new ZoneClimateReport(ZoneClimateStatus.adapt(item),
									 item.NumAux,
									 State.adapt(item.Current),
									 State.adapt(item.CurrentEnd),
									 State.adapt(item.Next),
									 State.adapt(item.NextEnd),
									 nextTime);

	}
}
