import { ClimateState } from './climate-state';
import testData from './unit-testdata/ClimateState.json';

describe('ClimateState', () => {
  it('should create an instance', () => {
    expect(new ClimateState()).toBeTruthy();
  });

  it('should be parsed from JSON', () => {
    for (const plain of testData) {
      let e = ClimateState.fromPlain(plain);
      expect(e.name).toEqual(plain.name || '');
      expect(e.temperature).toEqual(plain.temperature);
      expect(e.humidity).toEqual(plain.humidity);
      expect(e.wind).toEqual(plain.wind);
      expect(e.visible_light).toEqual(plain.visible_light);
      expect(e.uv_light).toEqual(plain.uv_light);
    }
  });
});
