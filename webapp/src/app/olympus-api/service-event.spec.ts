import { classToPlain, plainToClass } from 'class-transformer';
import { ServiceEvent, ServiceEventList, ServicesLogs } from './service-event';
import eventTestData from './unit-testdata/ServiceEvent.json';
import listTestData from './unit-testdata/ServiceEventList.json';
import logsTestData from './unit-testdata/ServicesLogs.json';

describe('ServiceEvent', () => {
  it('should be created', () => {
    expect(new ServiceEvent()).toBeTruthy();
  });

  it('should be parsed from JSON', () => {
    for (const plain of eventTestData) {
      let e = ServiceEvent.fromPlain(plain);
      expect(e).toBeTruthy();
      expect(e.time).toEqual(new Date(plain.time) || new Date(0));
      expect(e.on).toEqual(plain.on || false);
      expect(e.graceful).toEqual(plain.graceful || false);
    }
  });
});

describe('ServiceEventList', () => {
  it('should be created', () => {
    expect(new ServiceEventList()).toBeTruthy();
  });

  it('should be parsed from JSON', () => {
    for (const plain of listTestData) {
      let e = ServiceEventList.fromPlain(plain);
      expect(e).toBeTruthy();
      expect(e.zone).toEqual(plain.zone || '');
      expect(e.events).toEqual(
        (plain.events || []).map((v: any) => plainToClass(ServiceEvent, v))
      );
    }
  });
});

describe('ServicesLogs', () => {
  it('should be created', () => {
    expect(new ServicesLogs()).toBeTruthy();
  });

  it('should be parsed from JSON', () => {
    for (const plain of logsTestData) {
      let e = ServicesLogs.fromPlain(plain);
      expect(e).toBeTruthy();
      expect(e.climate).toEqual(
        (plain.climate || []).map((v: any) => plainToClass(ServiceEventList, v))
      );
      expect(e.tracking).toEqual(
        (plain.tracking || []).map((v: any) =>
          plainToClass(ServiceEventList, v)
        )
      );
    }
  });
});
