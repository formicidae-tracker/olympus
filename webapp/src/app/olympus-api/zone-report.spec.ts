import { ZoneReport } from './zone-report';
import testData from './unit-testdata/ZoneReport.json';
import { plainToClass } from 'class-transformer';
import { ZoneClimateReport } from './zone-climate-report';
import { TrackingInfo } from './tracking-info';
import { AlarmReport } from './alarm-report';

describe('ZoneReport', () => {
  it('should be created', () => {
    expect(new ZoneReport()).toBeTruthy();
  });

  it('should be parsed from JSON', () => {
    for (const plain of testData) {
      let e = ZoneReport.fromPlain(plain);
      expect(e).toBeTruthy();
      expect(e.host).toEqual(plain.host || '');
      expect(e.name).toEqual(plain.name || '');
      expect(e.climate).toEqual(plainToClass(ZoneClimateReport, plain.climate));
      expect(e.tracking).toEqual(plainToClass(TrackingInfo, plain.tracking));
      expect(e.alarms).toEqual(
        (plain.alarms || []).map((v: any) => plainToClass(AlarmReport, v))
      );
    }
  });
});
