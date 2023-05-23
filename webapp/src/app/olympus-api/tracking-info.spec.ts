import { cases } from 'jasmine-parameterized';

import { StreamInfo } from './stream-info';
import { TrackingInfo } from './tracking-info';
import testData from './unit-testdata/TrackingInfo.json';

describe('TrackingInfo', () => {
  it('should be created', () => {
    expect(new TrackingInfo()).toBeTruthy();
  });

  it('should be parsed from JSON', () => {
    for (const plain of testData) {
      let e = TrackingInfo.fromPlain(plain);
      expect(e).toBeTruthy();
      if (e == undefined) {
        continue;
      }
      expect(e.total_bytes).toEqual(plain.total_bytes || 0);
      expect(e.free_bytes).toEqual(plain.free_bytes || 0);
      expect(e.bytes_per_second).toEqual(plain.bytes_per_second || 0);
      expect(e.stream).toEqual(StreamInfo.fromPlain(plain.stream));
    }
  });

  cases([
    [{}, 0],
    [{ total_bytes: 1 }, 1],
    [{ total_bytes: 1, free_bytes: 1 }, 0],
    [{ total_bytes: 1, free_bytes: 2 }, 0],
    [{ free_bytes: 1 }, 0],
  ]).it('should provide filled_bytes accordingly', ([plain, expected]) => {
    expect(TrackingInfo.fromPlain(plain)?.used_bytes).toEqual(expected);
  });

  cases([
    [{ free_bytes: 10, bytes_per_second: 10 }, 1000],
    [{}, 0],
    [{ free_bytes: -1 }, 0],
    [{ free_bytes: 1 }, Infinity],
    [{ free_bytes: 1, bytes_per_second: -1 }, Infinity],
  ]).it('should compute filledUpETA', ([plain, expected]) => {
    expect(TrackingInfo.fromPlain(plain)?.filledUpEta()).toEqual(expected);
  });
});
