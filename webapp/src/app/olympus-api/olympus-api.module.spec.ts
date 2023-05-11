import {
  Bounds,
  ClimateState,
  StreamInfo,
  ZoneClimateReport,
  ZoneReportSummary,
} from './olympus-api.module';
import boundsData from './examples/Bounds.json';
import climateStateData from './examples/ClimateState.json';
import streamInfoData from './examples/StreamInfo.json';
import zoneClimateReportData from './examples/ZoneClimateReport.json';
import zoneReportSummaryData from './examples/ZoneReportSummary.json';

describe('Bounds', () => {
  it('should be parsable', () => {
    for (const raw of boundsData) {
      let parsed: Bounds = raw;
      expect(parsed).toEqual(raw);
    }
  });
});

describe('ClimateState', () => {
  it('should be parsable', () => {
    for (const raw of climateStateData) {
      let parsed: ClimateState = raw;
      expect(parsed).toEqual(raw);
    }
  });
});

describe('ZoneClimateReport', () => {
  it('should be parsable', () => {
    for (const raw of zoneClimateReportData) {
      let parsed: ZoneClimateReport = raw;
      expect(parsed).toEqual(raw);
    }
  });
});

describe('StreamInfo', () => {
  it('should be parsable', () => {
    for (const raw of streamInfoData) {
      let parsed: StreamInfo = raw;
      expect(parsed).toEqual(raw);
    }
  });
});

describe('ZoneReportSummary', () => {
  it('should be parsable', () => {
    for (const raw of zoneReportSummaryData) {
      let parsed: ZoneReportSummary = raw;
      expect(parsed).toEqual(raw);
    }
  });
});
