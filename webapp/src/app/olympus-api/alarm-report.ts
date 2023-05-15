import { Type, plainToClass } from 'class-transformer';
import { AlarmTimePoint } from './alarm-time-point';

export class AlarmReport {
  public identification: string = '';
  public level: number = 0;
  public description: string = '';

  @Type(() => AlarmTimePoint)
  public events: AlarmTimePoint[] = [];

  static fromPlain(plain: any): AlarmReport {
    return plainToClass(AlarmReport, plain, {
      exposeDefaultValues: true,
    });
  }
}
