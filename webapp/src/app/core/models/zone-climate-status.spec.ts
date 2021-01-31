import { ZoneClimateStatus } from './zone-climate-status';

describe('ZoneClimateStatus', () => {
	it('should create an instance', () => {
		expect(new ZoneClimateStatus(12,30)).toBeTruthy();
	});

	it('should adapt from different values',() => {
		expect(ZoneClimateStatus.adapt(null)).toBeNull();
		let s = ZoneClimateStatus.adapt({"Temperature":21.99,
										 "Humidity":60.36,
										 "TemperatureBounds":{
											 "Min":20,
											 "Max":null
										 },
										 "HumidityBounds":null,
										 "ActiveWarnings":12,
										 "ActiveEmergencies":35
										});
		expect(s.Temperature).toBe(21.99);
		expect(s.Humidity).toBe(60.36);
		expect(s.TemperatureBounds.Min).toBe(20);
		expect(s.TemperatureBounds.Max).toBe(100);
		expect(s.HumidityBounds.Min).toBe(0);
		expect(s.HumidityBounds.Max).toBe(100);
		expect(s.ActiveWarnings).toBe(12);
		expect(s.ActiveEmergencies).toBe(35);
	});

});