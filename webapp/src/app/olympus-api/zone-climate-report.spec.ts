import { ZoneClimateReport } from './zone-climate-report';
import { Bounds } from './bounds';
import { ClimateState } from './climate-state';

import testData from './unit-testdata/ZoneClimateReport.json';

describe('ZoneClimateReport', () => {
  it('should be created', () => {
    expect(new ZoneClimateReport()).toBeTruthy();
  });

  it('should be parsed from JSON', () => {
    for (const plain of testData) {
      let e = ZoneClimateReport.fromPlain(plain);
      expect(e).toBeTruthy();
      if (e == undefined) {
        continue;
      }
      expect(e.since).toEqual(new Date(plain.since || 0));
      expect(e.temperature).toEqual(plain.temperature);
      expect(e.humidity).toEqual(plain.humidity);

      expect(e.temperature_bounds).toEqual(
        Bounds.fromPlain(plain.temperature_bounds)
      );
      expect(e.humidity_bounds).toEqual(
        Bounds.fromPlain(plain.humidity_bounds)
      );

      expect(e.current).toEqual(ClimateState.fromPlain(plain.current));
      expect(e.current_end).toBeUndefined();

      expect(e.next).toEqual(ClimateState.fromPlain(plain.next));
      expect(e.next_end).toEqual(ClimateState.fromPlain(plain.next_end));

      if (plain.next_time != undefined) {
        expect(e.next_time).toEqual(new Date(plain.next_time));
      } else {
        expect(e.next_time).toBeUndefined();
      }
    }
  });
});
