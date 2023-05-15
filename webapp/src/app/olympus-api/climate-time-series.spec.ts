import { Point, ClimateTimeSeries } from './climate-time-series';
import testData from './unit-testdata/ClimateTimeSeries.json';

describe('ClimateTimeSeries', () => {
  it('should be created', () => {
    expect(new ClimateTimeSeries()).toBeTruthy();
  });

  it('should be parsed from JSON', () => {
    for (const plain of testData) {
      let e = ClimateTimeSeries.fromPlain(plain);
      expect(e).toBeTruthy();
      expect(e.units).toEqual(plain.units || '');
      expect(e.reference).toEqual(new Date(plain.reference) || new Date(0));
      expect(e.humidity).toEqual(
        (plain.humidity || []).map(
          (v: any) => new Point(v.x || 0.0, v.y || 0.0)
        )
      );
      expect(e.temperature).toEqual(
        (plain.temperature || []).map(
          (v: any) => new Point(v.x || 0.0, v.y || 0.0)
        )
      );
    }
  });
});
