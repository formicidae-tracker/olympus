import { ZoneClimateStatus } from './zone-climate-status';

describe('ZoneClimateStatus', () => {
	it('should create an instance', () => {
		expect(new ZoneClimateStatus(12,30)).toBeTruthy();
	});

	it('should adapt from different values',() => {
		expect(ZoneClimateStatus.adapt(null)).toBeNull();
		let s = ZoneClimateStatus.adapt(
			{"Temperature":21.99,
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
		expect(s.TemperatureBounds.Max).toBeNaN();
		expect(s.HumidityBounds.Min).toBeNaN();
		expect(s.HumidityBounds.Max).toBeNaN();
		expect(s.ActiveWarnings).toBe(12);
		expect(s.ActiveEmergencies).toBe(35);
		s = ZoneClimateStatus.adapt({
			"Temperature":-1000,
			"Humidity":-1000,
			"TemperatureBounds":{
				"Min":null,
				"Max":null
			},
			"HumidityBounds":{
				"Min":null,
				"Max":null
			},
			"ActiveWarnings":0,
			"ActiveEmergencies":1
		})
		expect(s.Temperature).toBeNaN()
		expect(s.Humidity).toBeNaN()
		expect(s.HumidityBounds.Min).toBeNaN();
		expect(s.HumidityBounds.Max).toBeNaN();
		expect(s.TemperatureBounds.Min).toBeNaN();
		expect(s.TemperatureBounds.Max).toBeNaN();
		expect(s.ActiveWarnings).toBe(0);
		expect(s.ActiveEmergencies).toBe(1);


	});

});
