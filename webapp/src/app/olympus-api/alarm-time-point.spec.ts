import { classToPlain, plainToClass } from 'class-transformer';
import { AlarmTimePoint } from './alarm-time-point';
import testData from './unit-testdata/AlarmTimePoint.json';

describe('AlarmTimePoint', () => {
  it('should be created', () => {
    expect(new AlarmTimePoint()).toBeTruthy();
  });

  it('should be parsed from JSON', () => {
    for (const plain of testData) {
      let e = AlarmTimePoint.fromPlain(plain);
      expect(e).toBeTruthy();
      expect(e.time).toEqual(new Date(plain.time) || new Date(0));
      expect(e.on).toEqual(plain.on || false);
    }
  });
});
