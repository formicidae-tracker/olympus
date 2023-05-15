import { AlarmReport } from './alarm-report';
import { TrackingInfo } from './tracking-info';
import { ZoneClimateReport } from './zone-climate-report';

export class ZoneReport {
  public host: string = '';
  public name: string = '';
  public climate?: ZoneClimateReport;
  public tracking?: TrackingInfo;
  public alarms: AlarmReport[] = [];

  static fromPlain(plain: any): ZoneReport {
    let res = new ZoneReport();
    res.host = plain.host || '';
    res.name = plain.name || '';
    if (plain.climate != undefined) {
      res.climate = ZoneClimateReport.fromPlain(plain.climate);
    }
    if (plain.tracking != undefined) {
      res.tracking = TrackingInfo.fromPlain(plain.tracking);
    }
    for (const a of plain.alarms || []) {
      res.alarms.push(AlarmReport.fromPlain(a));
    }
    return res;
  }
}
