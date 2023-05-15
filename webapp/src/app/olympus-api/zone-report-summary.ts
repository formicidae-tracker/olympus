import { Type, plainToClass } from 'class-transformer';
import { ZoneClimateReport } from './zone-climate-report';
import { TrackingInfo } from './tracking-info';

export class ZoneReportSummary {
  public host: string = '';
  public name: string = '';

  @Type(() => ZoneClimateReport)
  climate?: ZoneClimateReport;
  @Type(() => TrackingInfo)
  tracking?: TrackingInfo;

  active_warnings: number = 0;
  active_emergencies: number = 0;

  static fromPlain(plain: any): ZoneReportSummary {
    return plainToClass(ZoneReportSummary, plain, {
      exposeDefaultValues: true,
    });
  }
}
