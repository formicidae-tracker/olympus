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

  it('can be sorted', () => {
    let a = AlarmReport.fromPlain({
      level: 1,
      events: [{ time: new Date(1), on: true }],
    });
    let b = AlarmReport.fromPlain({
      level: 1,
      events: [{ time: new Date(0), on: true }],
    });
    let c = AlarmReport.fromPlain({
      level: 0,
      events: [{ time: new Date(1), on: true }],
    });
    let d = AlarmReport.fromPlain({
      level: 0,
      events: [{ time: new Date(0), on: true }],
    });
    let e = AlarmReport.fromPlain({
      level: 1,
      events: [{ time: new Date(1), on: false }],
    });
    let f = AlarmReport.fromPlain({
      level: 1,
      events: [{ time: new Date(0), on: false }],
    });
    let g = AlarmReport.fromPlain({
      level: 0,
      events: [{ time: new Date(1), on: false }],
    });
    let h = AlarmReport.fromPlain({
      level: 0,
      events: [{ time: new Date(0), on: false }],
    });

    let sorted = [a, b, c, d, e, f, g, h];
    let unsorted = [b, e, h, a, c, d, g, f];
    expect(unsorted.sort(AlarmReport.compareFunction)).toEqual(sorted);
  });
});
