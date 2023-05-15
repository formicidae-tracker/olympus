import { ZoneReportSummary } from './zone-report-summary';
import { ZoneClimateReport } from './zone-climate-report';
import { TrackingInfo } from './tracking-info';

import testData from './unit-testdata/ZoneReportSummary.json';

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
      expect(e.climate).toEqual(ZoneClimateReport.fromPlain(plain.climate));
      expect(e.tracking).toEqual(TrackingInfo.fromPlain(plain.tracking));
      expect(e.active_warnings).toEqual(plain.active_warnings || 0);
      expect(e.active_emergencies).toEqual(plain.active_emergencies || 0);
    }
  });
});
