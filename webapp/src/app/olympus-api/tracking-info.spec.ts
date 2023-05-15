import { plainToClass } from 'class-transformer';
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
      expect(e.total_bytes).toEqual(plain.total_bytes || 0);
      expect(e.free_bytes).toEqual(plain.free_bytes || 0);
      expect(e.bytes_per_second).toEqual(plain.bytes_per_second || 0);
      expect(e.stream).toEqual(
        plainToClass(StreamInfo, plain.stream) || undefined
      );
    }
  });
});
