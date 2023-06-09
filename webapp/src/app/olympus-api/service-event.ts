import { Event } from './event';

export class ServiceLog {
  public zone: string = '';
  public events: Event[] = [];

  static fromPlain(plain: any): ServiceLog {
    let res = new ServiceLog();
    res.zone = plain.zone || '';
    for (const e of plain.events || []) {
      res.events.push(Event.fromPlain(e));
    }
    return res;
  }
}
