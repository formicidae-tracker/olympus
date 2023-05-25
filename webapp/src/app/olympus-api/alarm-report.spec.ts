import { cases } from 'jasmine-parameterized';

import { AlarmEvent, AlarmReport } from './alarm-report';

import testData from './unit-testdata/AlarmReport.json';

describe('AlarmEvent', () => {
  cases([
    [new Date(0), undefined],
    [new Date(0), new Date(1)],
  ]).it('should create', ([start, end]) => {
    expect(new AlarmEvent(start, end)).toBeTruthy();
  });

  cases([
    [new AlarmEvent(new Date(0), undefined), true],
    [new AlarmEvent(new Date(0), new Date(1)), false],
  ]).it('should be on when it has an end Date', ([event, expected]) => {
    expect(event.on()).toEqual(expected);
  });

  cases([
    [new AlarmEvent(new Date(0), undefined), new Date(0)],
    [new AlarmEvent(new Date(0), new Date(1)), new Date(1)],
  ]).it('should be on when it has an end Date', ([event, expected]) => {
    expect(event.time()).toEqual(expected);
  });
});

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
    }
  });

  it('can be sorted', () => {
    let a = AlarmReport.fromPlain({
      level: 1,
      events: [{ start: new Date(1) }],
    });
    console.log(a);
    let b = AlarmReport.fromPlain({
      level: 1,
      events: [{ start: new Date(0) }],
    });
    let c = AlarmReport.fromPlain({
      level: 0,
      events: [{ start: new Date(1) }],
    });
    let d = AlarmReport.fromPlain({
      level: 0,
      events: [{ start: new Date(0) }],
    });
    let e = AlarmReport.fromPlain({
      level: 1,
      events: [{ start: new Date(3), end: new Date(4) }],
    });
    let f = AlarmReport.fromPlain({
      level: 1,
      events: [{ start: new Date(2), end: new Date(3) }],
    });
    let g = AlarmReport.fromPlain({
      level: 0,
      events: [{ start: new Date(2), end: new Date(4) }],
    });
    let h = AlarmReport.fromPlain({
      level: 0,
      events: [{ start: new Date(1), end: new Date(3) }],
    });

    let sorted = [a, b, c, d, e, f, g, h];
    let unsorted = [b, e, h, a, c, d, g, f];
    expect(unsorted.sort(AlarmReport.compareFunction)).toEqual(sorted);
  });
});
