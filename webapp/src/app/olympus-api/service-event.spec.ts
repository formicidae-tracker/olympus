import { Event } from './event';
import { ServiceEventList, ServicesLogs } from './service-event';

import eventTestData from './unit-testdata/ServiceEvent.json';
import listTestData from './unit-testdata/ServiceEventList.json';
import logsTestData from './unit-testdata/ServicesLogs.json';

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
        (plain.events || []).map((v: any) => Event.fromPlain(v))
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
        (plain.climate || []).map((v: any) => ServiceEventList.fromPlain(v))
      );
      expect(e.tracking).toEqual(
        (plain.tracking || []).map((v: any) => ServiceEventList.fromPlain(v))
      );
    }
  });
});
