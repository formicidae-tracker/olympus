import { ZoneSummaryReport } from './zone-summary-report';

describe('ZoneSummaryReport', () => {
  it('should create an instance', () => {
      expect(new ZoneSummaryReport('foo',
								   'bar')).toBeTruthy();
  });

});
