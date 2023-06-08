export class ServiceEvent {
  public start: Date = new Date(0);
  public end?: Date;
  public graceful: boolean = false;

  static fromPlain(plain: any): ServiceEvent {
    let res = new ServiceEvent();
    res.start = new Date(plain.start || 0);
    if (plain.end != undefined) {
      res.end = new Date(plain.end);
    }
    res.graceful = plain.graceful || false;
    return res;
  }
}

export class ServiceEventList {
  public zone: string = '';
  public events: ServiceEvent[] = [];

  static fromPlain(plain: any): ServiceEventList {
    let res = new ServiceEventList();
    res.zone = plain.zone || '';
    for (const e of plain.events || []) {
      res.events.push(ServiceEvent.fromPlain(e));
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
