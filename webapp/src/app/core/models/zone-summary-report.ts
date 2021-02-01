import { StreamInfo } from './stream-info';
import { ZoneClimateReport } from './zone-climate-report';

export class ZoneSummaryReport {
	constructor(public host: string,
				public zoneName: string,
				public streamInfo: StreamInfo  = new StreamInfo(),
				public climate: ZoneClimateReport = null ) {
	}

	static adapt(item: any): ZoneSummaryReport {
		return new ZoneSummaryReport(item.Host,
									 item.Name,
									 StreamInfo.adapt(item.StreamURL),
									 ZoneClimateReport.adapt(item.Climate));
	}

	hasStream(): boolean {
		return this.streamInfo.hasStream()
	}

	hasThumbnail(): boolean {
		return this.streamInfo.hasThumbnail()
	}

	streamURL(): string {
		return this.streamInfo.streamURL;
	}

	thumbnailURL(): string {
		return this.streamInfo.thumbnailURL;
	}

	hasClimate(): boolean {
		return this.climate != null;
	}

}
