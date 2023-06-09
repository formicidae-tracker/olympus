import { cases } from 'jasmine-parameterized';

import { Event } from './event';

describe('Event', () => {
  cases([
    [new Date(0), undefined, undefined],
    [new Date(0), new Date(1), undefined],
    [new Date(0), new Date(1), true],
    [new Date(0), new Date(1), false],
  ]).it('should create', ([start, end, graceful]) => {
    expect(new Event(start, end, graceful)).toBeTruthy();
  });

  cases([
    [new Event(new Date(0), undefined), true, 'running'],
    [new Event(new Date(0), new Date(1)), false, 'done'],
    [new Event(new Date(0), new Date(1), false), false, 'ungraceful'],
    [new Event(new Date(0), new Date(1), true), false, 'done'],
  ]).it(
    'should be on when it has an end Date',
    ([event, expectedOn, expectedStatus]) => {
      expect(event.on()).toEqual(expectedOn);
      expect(event.status()).toEqual(expectedStatus);
    }
  );

  cases([
    [new Event(new Date(0), undefined), new Date(0)],
    [new Event(new Date(0), new Date(1)), new Date(1)],
  ]).it('should report meaningful time', ([event, expected]) => {
    expect(event.time()).toEqual(expected);
  });

  cases([
    [new Event(new Date(0), undefined), undefined],
    [new Event(new Date(0), new Date(1)), 1],
  ]).it('should report duration', ([event, expected]) => {
    expect(event.duration()).toEqual(expected);
  });

  cases([
    ['{"start":"2023-03-01T00:00:00.000Z"}'],
    ['{"start":"2023-03-01T00:00:00.000Z","end":"2023-03-01T01:00:00.000Z"}'],
    [
      '{"start":"2023-03-01T00:00:00.000Z","end":"2023-03-01T01:00:00.000Z","graceful":true}',
    ],
    [
      '{"start":"2023-03-01T00:00:00.000Z","end":"2023-03-01T01:00:00.000Z","graceful":false}',
    ],
  ]).it('should parse from JSON', ([jsondata]) => {
    const plain = JSON.parse(jsondata);
    const e = Event.fromPlain(plain);
    expect(e).toBeTruthy();
    expect(e.start).toEqual(new Date(plain.start || 0));
    if (plain.end != undefined) {
      expect(e.end).toEqual(new Date(plain.end));
    }
    expect(e.graceful).toEqual(plain.graceful);
  });
});
