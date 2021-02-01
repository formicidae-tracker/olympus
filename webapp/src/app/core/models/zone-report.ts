import { StreamInfo } from "./stream-info";
import { ZoneClimateReport } from "./zone-climate-report";

export class ZoneReport {
	constructor(public host: string,
				public name: string,
				public climate: ZoneClimateReport = null,
				public streamInfo: StreamInfo = null,
				public alarms: any[] = null) {
	}

	static adapt(item: any) {
		return new ZoneReport(item.Host,
							  item.Name,
							  ZoneClimateReport.adapt(item.Climate),
							  StreamInfo.adapt(item.StreamInfo),
							  item.Alarms);
	}
}
