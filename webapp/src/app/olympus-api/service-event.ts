import { Event } from './event';

export class ServiceEventList {
  public zone: string = '';
  public events: Event[] = [];

  static fromPlain(plain: any): ServiceEventList {
    let res = new ServiceEventList();
    res.zone = plain.zone || '';
    for (const e of plain.events || []) {
      res.events.push(Event.fromPlain(e));
    }
    return res;
  }
}

export class ServicesLogs {
  public climate: ServiceEventList[] = [];
  public tracking: ServiceEventList[] = [];

  static fromPlain(plain: any): ServicesLogs {
    let res = new ServicesLogs();
    for (const l of plain.climate || []) {
      res.climate.push(ServiceEventList.fromPlain(l));
    }
    for (const l of plain.tracking || []) {
      res.tracking.push(ServiceEventList.fromPlain(l));
    }
    return res;
  }
}
