import { StreamInfo } from './stream-info';
import testData from './unit-testdata/StreamInfo.json';

describe('StreamInfo', () => {
  it('should be created', () => {
    expect(new StreamInfo()).toBeTruthy();
  });

  it('should be parsed from JSON', () => {
    for (const plain of testData) {
      let e = StreamInfo.fromPlain(plain);
      expect(e).toBeTruthy();
      if (e != undefined) {
        expect(e.experiment_name).toEqual(plain.experiment_name || '');
        expect(e.stream_URL).toEqual(plain.stream_URL || '');
        expect(e.thumbnail_URL).toEqual(plain.thumbnail_URL || '');
      }
    }
  });
});
