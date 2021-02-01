
import { ZoneClimateStatus } from './zone-climate-status';
import { State } from './state';
import { StreamInfo } from './stream-info';

export class ZoneReport {
	constructor(public ClimateStatus: ZoneClimateStatus,
				public StreamInfo: StreamInfo,
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
