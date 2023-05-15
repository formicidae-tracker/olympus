import { ZoneClimateReport } from './zone-climate-report';
import { plainToClass } from 'class-transformer';
import { Bounds } from './bounds';

import testData from './unit-testdata/ZoneClimateReport.json';
import { ClimateState } from './climate-state';

describe('ZoneClimateReport', () => {
  it('should be created', () => {
    expect(new ZoneClimateReport()).toBeTruthy();
  });

  it('should be parsed from JSON', () => {
    for (const plain of testData) {
      let e = ZoneClimateReport.fromPlain(plain);
      expect(e).toBeTruthy();
      expect(e.temperature).toEqual(plain.temperature);
      expect(e.humidity).toEqual(plain.humidity);

      expect(e.temperature_bounds).toEqual(
        plainToClass(Bounds, plain.temperature_bounds)
      );
      expect(e.humidity_bounds).toEqual(
        plainToClass(Bounds, plain.humidity_bounds)
      );

      expect(e.current).toEqual(plainToClass(ClimateState, plain.current));
      expect(e.current_end).toBeUndefined();

      expect(e.next).toEqual(plainToClass(ClimateState, plain.next));
      expect(e.next_end).toEqual(plainToClass(ClimateState, plain.next_end));
      if (plain.next_time != undefined) {
        expect(e.next_time).toEqual(new Date(plain.next_time));
      } else {
        expect(e.next_time).toBeUndefined();
      }
    }
  });
});
