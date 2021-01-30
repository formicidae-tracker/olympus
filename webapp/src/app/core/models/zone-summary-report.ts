import { Injectable } from '@angular/core';
import { Adapter } from './adapter';

import { ZoneClimateStatus,ZoneClimateStatusAdapter } from './zone-climate-status';


export class ZoneSummaryReport {
	constructor(public Host: string,
				public ZoneName: string,
				public StreamURL: string  = '',
				public Status: ZoneClimateStatus = null ) {
	}
}

@Injectable({
    providedIn: 'root'
})

export class ZoneSummaryReportAdapter implements Adapter<ZoneSummaryReport> {
	constructor(private zcsAdapter :ZoneClimateStatusAdapter){}

	adapt(item: any): ZoneSummaryReport {
		return new ZoneSummaryReport(item.Host,
									 item.Name,
									 item.StreamURL,
									 this.zcsAdapter.adapt(item.Climate));

	}
}
