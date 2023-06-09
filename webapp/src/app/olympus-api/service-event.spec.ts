import { Event } from './event';
import { ServiceLog } from './service-event';

import testData from './unit-testdata/ServiceLog.json';

describe('ServiceLog', () => {
  it('should be created', () => {
    expect(new ServiceLog()).toBeTruthy();
  });

  it('should be parsed from JSON', () => {
    for (const plain of testData) {
      let e = ServiceLog.fromPlain(plain);
      expect(e).toBeTruthy();
      expect(e.zone).toEqual(plain.zone || '');
      expect(e.events).toEqual(
        (plain.events || []).map((v: any) => Event.fromPlain(v))
      );
    }
  });
});
