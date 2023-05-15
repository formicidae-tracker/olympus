import { Type, plainToClass } from 'class-transformer';
import { AlarmReport } from './alarm-report';
import { TrackingInfo } from './tracking-info';
import { ZoneClimateReport } from './zone-climate-report';

export class ZoneReport {
  public host: string = '';
  public name: string = '';
  @Type(() => ZoneClimateReport)
  public climate?: ZoneClimateReport;
  @Type(() => TrackingInfo)
  public tracking?: TrackingInfo;
  @Type(() => AlarmReport)
  public alarms: AlarmReport[] = [];

  static fromPlain(plain: any): ZoneReport {
    return plainToClass(ZoneReport, plain, {
      exposeDefaultValues: true,
    });
  }
}
