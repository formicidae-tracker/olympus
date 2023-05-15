import { Bounds } from './bounds';
import testData from './unit-testdata/Bounds.json';

describe('Bounds', () => {
  it('should create an instance', () => {
    expect(new Bounds()).toBeTruthy();
  });

  it('should be parsed from JSON', () => {
    for (const plain of testData) {
      let b = Bounds.fromPlain(plain);
      expect(b).toBeTruthy();
      expect(b.minimum).toEqual(plain.minimum);
      expect(b.maximum).toEqual(plain.maximum);
    }
  });
});
