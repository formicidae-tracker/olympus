import { Type, plainToClass } from 'class-transformer';

export class ServiceEvent {
  @Type(() => Date)
  public time: Date = new Date(0);
  public on: boolean = false;
  public graceful: boolean = false;

  static fromPlain(plain: any): ServiceEvent {
    return plainToClass(ServiceEvent, plain, {
      exposeDefaultValues: true,
    });
  }
}

export class ServiceEventList {
  public zone: string = '';
  @Type(() => ServiceEvent)
  public events: ServiceEvent[] = [];

  static fromPlain(plain: any): ServiceEventList {
    return plainToClass(ServiceEventList, plain, {
      exposeDefaultValues: true,
    });
  }
}

export class ServicesLogs {
  @Type(() => ServiceEventList)
  public climate: ServiceEventList[] = [];
  @Type(() => ServiceEventList)
  public tracking: ServiceEventList[] = [];

  static fromPlain(plain: any): ServicesLogs {
    return plainToClass(ServicesLogs, plain, {
      exposeDefaultValues: true,
    });
  }
}
