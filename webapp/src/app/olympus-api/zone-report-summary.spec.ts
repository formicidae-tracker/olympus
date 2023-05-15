import { ZoneReportSummary } from './zone-report-summary';
import testData from './unit-testdata/ZoneReportSummary.json';
import { plainToClass } from 'class-transformer';
import { ZoneClimateReport } from './zone-climate-report';
import { TrackingInfo } from './tracking-info';

describe('ZoneReportSummary', () => {
  it('should be created', () => {
    expect(new ZoneReportSummary()).toBeTruthy();
  });

  it('should be parsed from JSON', () => {
    for (const plain of testData) {
      let e = ZoneReportSummary.fromPlain(plain);
      expect(e).toBeTruthy();
      expect(e.host).toEqual(plain.host || '');
      expect(e.name).toEqual(plain.name || '');
      expect(e.climate).toEqual(plainToClass(ZoneClimateReport, plain.climate));
      expect(e.tracking).toEqual(plainToClass(TrackingInfo, plain.tracking));
      expect(e.active_warnings).toEqual(plain.active_warnings || 0);
      expect(e.active_emergencies).toEqual(plain.active_emergencies || 0);
    }
  });
});
