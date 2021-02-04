import { AlarmReport } from "./alarm";
import { StreamInfo } from "./stream-info";
import { ZoneClimateReport } from "./zone-climate-report";

export class ZoneReport {
	constructor(public host: string,
				public name: string,
				public climate: ZoneClimateReport = null,
				public streamInfo: StreamInfo = null,
				public alarms: AlarmReport[] = []) {
	}

	static adapt(item: any) {
		let alarms: AlarmReport[] = [];
		if ( item.Alarms != null ) {
			for ( let a of item.Alarms ) {
				alarms.push(AlarmReport.adapt(a));
			}
		}

		return new ZoneReport(item.Host,
							  item.Name,
							  ZoneClimateReport.adapt(item.Climate),
							  StreamInfo.adapt(item.Stream),
							  alarms);
	}
}
