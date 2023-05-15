import { AlarmReport } from './alarm-report';
import testData from './unit-testdata/AlarmReport.json';

describe('AlarmReport', () => {
  it('should be created', () => {
    expect(new AlarmReport()).toBeTruthy();
  });

  it('should be parsed from JSON', () => {
    for (const plain of testData) {
      let e = AlarmReport.fromPlain(plain);
      expect(e).toBeTruthy();
      expect(e.identification).toEqual(plain.identification || '');
      expect(e.level).toEqual(plain.level || 0);
      expect(e.description).toEqual(plain.description || '');
      expect(e.events).toBeDefined();
    }
  });
});
