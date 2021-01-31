import { ZoneClimateStatus } from './zone-climate-status';


export class ZoneSummaryReport {
	constructor(public Host: string,
				public ZoneName: string,
				public StreamURL: string  = '',
				public Status: ZoneClimateStatus = null ) {
	}

	static adapt(item: any): ZoneSummaryReport {
		return new ZoneSummaryReport(item.Host,
									 item.Name,
									 item.StreamURL,
									 ZoneClimateStatus.adapt(item.Climate));
	}
}
