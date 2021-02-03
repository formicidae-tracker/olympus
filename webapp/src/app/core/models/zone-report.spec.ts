import { StreamInfo } from './stream-info';
import { ZoneReport } from './zone-report';



describe('ZoneReport', () => {

	it('should create a new instance', () => {
		expect(new ZoneReport('foo', 'bar')).toBeTruthy()
	});


	let jsonData = [
		{
			"Host": "atreides",
			"Name": "box",
			"Climate": {
				"Temperature": 26.00,
				"Humidity": 61.48,
				"TemperatureBounds": {
					"Min": 20,
					"Max": 31
				},
				"HumidityBounds": {
					"Min": 40,
					"Max": 80
				},
				"ActiveWarnings": 0,
				"ActiveEmergencies": 0,
				"NumAux": 0,
				"Current": { "Name": "day", "Temperature": 26, "Humidity": 60, "Wind": 100, "VisibleLight": 40, "UVLight": 100 },
				"CurrentEnd": null, "Next": { "Name": "day to night", "Temperature": 26, "Humidity": 60, "Wind": 100, "VisibleLight": 40, "UVLight": 100 },
				"NextEnd": { "Name": "night", "Temperature": 22, "Humidity": 60, "Wind": 100, "VisibleLight": 0, "UVLight": 0 },
				"NextTime": "2021-02-05T17:00:00Z"
			},
			"Stream": null,
			"Alarms": [
				{
					"Reason": "Cannot reach desired humidity",
					"Level": 1,
					"Events": [
						{ "On": true, "Time": "2021-02-03T16:43:28.659707308+01:00" },
						{ "On": false, "Time": "2021-02-05T14:11:41.659707308+01:00" }
					]
				},
				{
					"Reason": "Celaeno is empty",
					"Level": 2, "Events": [
						{ "On": true, "Time": "2021-02-03T16:45:28.659707308+01:00" },
						{ "On": false, "Time": "2021-02-05T14:25:22.659707308+01:00" }
					]
				}
			]
		},
		{
			"Host": "fremens",
			"Name": "box",
			"Climate": null,
			"Stream": {
				"StreamURL": "/olympus/hls/fremens.m3u8",
				"ThumbnailURL": "/olympus/fremens.png"
			},
			"Alarms": null,
		},
	]

	describe('atreides', () => {
		let atreides: ZoneReport;
		beforeEach(() => {
			atreides = ZoneReport.adapt(jsonData[0]);
		});

		it('should be truthy', () => {
			expect(atreides).toBeTruthy();
		});

		it('should have empty stream info',() => {
			expect(atreides.streamInfo).toEqual(new StreamInfo());
		});

		it('should have a climate',() => {
			expect(atreides).toBeTruthy();
		});

		it('should have alarm reports',() => {
			expect(atreides.alarms.length).toBe(2,'alarms');
			if (atreides.alarms.length < 2) {
				return;
			}
			expect(atreides.alarms[0].reason).toContain('Cannot reach desired humidity');
			expect(atreides.alarms[1].reason).toContain('Celaeno is empty');
		});

	});

	describe('fremens', () => {
		let fremens: ZoneReport;
		beforeEach(() => {
			fremens = ZoneReport.adapt(jsonData[1]);
		});

		it('should be truthy', () => {
			expect(fremens).toBeTruthy();
		});
	});


});
