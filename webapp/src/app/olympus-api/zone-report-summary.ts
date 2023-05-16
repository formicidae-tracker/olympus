import { ZoneClimateReport } from './zone-climate-report';
import { TrackingInfo } from './tracking-info';

export class ZoneReportSummary {
  public host: string = '';
  public name: string = '';

  climate?: ZoneClimateReport;
  tracking?: TrackingInfo;

  active_warnings: number = 0;
  active_emergencies: number = 0;

  public identifier(): string {
    return this.host + '.' + this.name;
  }

  public streamThumbnailURL(): string | undefined {
    if (this.tracking == undefined || this.tracking.stream == undefined) {
      return undefined;
    }
    return this.tracking.stream.thumbnail_URL;
  }

  static fromPlain(plain: any): ZoneReportSummary {
    let ret = new ZoneReportSummary();
    ret.host = plain.host || '';
    ret.name = plain.name || '';
    if (plain.climate != undefined) {
      ret.climate = ZoneClimateReport.fromPlain(plain.climate);
    }
    if (plain.tracking != undefined) {
      ret.tracking = TrackingInfo.fromPlain(plain.tracking);
    }
    ret.active_emergencies = plain.active_emergencies || 0;
    ret.active_warnings = plain.active_warnings || 0;
    return ret;
  }
}
