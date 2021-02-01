import { ZoneClimateReport } from './zone-climate-report';
import { ZoneSummaryReport } from './zone-summary-report';

describe('ZoneSummaryReport', () => {
	it('should create an instance', () => {
		expect(new ZoneSummaryReport('foo',
									 'bar')).toBeTruthy();
	});

	it('should adapt from JSON', () => {
		let jsonData = [
			{
				"Host":"atreides",
				"Name":"box",
				"Climate":{
					"Temperature":26.00,
					"Humidity":59.60,
					"TemperatureBounds":{
						"Min":20,
						"Max":31
					},
					"HumidityBounds":{
						"Min":40,
						"Max":80},
					"ActiveWarnings":1,
					"ActiveEmergencies":1},
				"Stream":null
			},
			{
				"Host":"fremens",
				"Name":"box",
				"Climate":{
					"Temperature":26.01,
					"Humidity":60.78,
					"TemperatureBounds":{
						"Min":20,
						"Max":31
					},
					"HumidityBounds":{
						"Min":40,
						"Max":80
					},
					"ActiveWarnings":1,
					"ActiveEmergencies":1
				},
				"Stream":null
			},
			{
				"Host":"fremens",
				"Name":"tunnel",
				"Climate":{
					"Temperature":-1000.00,
					"Humidity":-1000.00,
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
				},
				"Stream":null
			}
		]

		for ( let d of jsonData ) {
			expect(ZoneClimateReport.adapt(d)).toBeTruthy();
		}

	});

});
