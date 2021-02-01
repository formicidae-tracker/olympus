import { ZoneClimateStatus } from './zone-climate-status';
import { StreamInfo } from './stream-info';

export class ZoneSummaryReport {
	constructor(public host: string,
				public zoneName: string,
				public streamURL: StreamInfo  = null,
				public status: ZoneClimateStatus = null ) {
	}

	static adapt(item: any): ZoneSummaryReport {
		return new ZoneSummaryReport(item.Host,
									 item.Name,
									 item.StreamURL,
									 ZoneClimateStatus.adapt(item.Climate));
	}
}
