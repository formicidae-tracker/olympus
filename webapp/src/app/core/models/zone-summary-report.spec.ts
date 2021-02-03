import { ZoneClimateReport } from './zone-climate-report';
import { ZoneSummaryReport } from './zone-summary-report';

describe('ZoneSummaryReport', () => {
	it('should create an instance', () => {
		expect(new ZoneSummaryReport('foo',
			'bar')).toBeTruthy();
	});

	let jsonData = [
		{
			"Host": "atreides",
			"Name": "box",
			"Climate": {
				"Temperature": 26.00,
				"Humidity": 59.60,
				"TemperatureBounds": {
					"Min": 20,
					"Max": 31
				},
				"HumidityBounds": {
					"Min": 40,
					"Max": 80
				},
				"ActiveWarnings": 1,
				"ActiveEmergencies": 1
			},
			"Stream": { "StreamURL": "/olympus/hls/atreides.m3u8", "ThumbnailURL": "/olympus/atreides.png" },
		},
		{
			"Host": "fremens",
			"Name": "box",
			"Climate": {
				"Temperature": 26.01,
				"Humidity": 60.78,
				"TemperatureBounds": {
					"Min": 20,
					"Max": 31
				},
				"HumidityBounds": {
					"Min": 40,
					"Max": 80
				},
				"ActiveWarnings": 1,
				"ActiveEmergencies": 1
			},
			"Stream": null
		},
		{
			"Host": "fremens",
			"Name": "tunnel",
			"Climate": {
				"Temperature": -1000.00,
				"Humidity": -1000.00,
				"TemperatureBounds": {
					"Min": null,
					"Max": null
				},
				"HumidityBounds": {
					"Min": null,
					"Max": null
				},
				"ActiveWarnings": 0,
				"ActiveEmergencies": 1
			},
			"Stream": null
		}
	]


	describe('atreides', () => {
		let atreides: ZoneSummaryReport;
		beforeEach(() => {
			atreides = ZoneSummaryReport.adapt(jsonData[0]);
		});

		it('should be truthy', () => {
			expect(atreides).toBeTruthy();
		});

		it('should have empty stream info',() => {
			expect(atreides.streamInfo.streamURL).toEqual('/olympus/hls/atreides.m3u8','streamURL');
			expect(atreides.streamInfo.thumbnailURL).toEqual('/olympus/atreides.png','thumbnailURL');
		});

		it('should have a climate',() => {
			expect(atreides).toBeTruthy();
		});


	});

	describe('fremens', () => {
		let fremens: ZoneSummaryReport;
		beforeEach(() => {
			fremens = ZoneSummaryReport.adapt(jsonData[1]);
		});

		it('should be truthy', () => {
			expect(fremens).toBeTruthy();
		});
	});

});
