import { Injectable } from '@angular/core';
import { Adapter } from './adapter';

import { ZoneClimateStatus,ZoneClimateStatusAdapter } from './zone-climate-status';
import { State,StateAdapter } from './state';

export class ZoneClimateReport {
	constructor(public ClimateStatus: ZoneClimateStatus,
				public NumAux: number,
				public Current: State,
				public CurrentEnd: State,
				public Next: State,
				public NextEnd: State,
				public NextTime: Date){}
}

Injectable({
	providedIn: 'root'
})

export class ZoneClimateReportAdapter implements Adapter<ZoneClimateReport> {
	constructor(private zcsAdapter: ZoneClimateStatusAdapter,
				private stateAdapter: StateAdapter){}

	adapt(item: any): ZoneClimateReport {
		let nextTime: Date = null;
		if ( item.NextTime != null ) {
			nextTime = new Date(item.NextTime);
		}
		return new ZoneClimateReport(this.zcsAdapter.adapt(item),
									 item.NumAux,
									 this.stateAdapter.adapt(item.Current),
									 this.stateAdapter.adapt(item.CurrentEnd),
									 this.stateAdapter.adapt(item.Next),
									 this.stateAdapter.adapt(item.NextEnd),
									 nextTime);
	}
}
