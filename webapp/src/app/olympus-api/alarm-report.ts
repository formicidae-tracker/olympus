import { AlarmTimePoint } from './alarm-time-point';

export class AlarmReport {
  public identification: string = '';
  public level: number = 0;
  public description: string = '';

  public events: AlarmTimePoint[] = [];

  static fromPlain(plain: any): AlarmReport {
    let res = new AlarmReport();
    res.identification = plain.identification || '';
    res.level = plain.level || 0;
    res.description = plain.description || '';
    for (const e of plain.events || []) {
      res.events.push(AlarmTimePoint.fromPlain(e));
    }
    return res;
  }
}
