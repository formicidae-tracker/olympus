import { ZoneReport } from './zone-report';



describe('ZoneReport', () => {

	it('should create a new instance', () => {
		expect(new ZoneReport('foo','bar')).toBeTruthy()
	});

	it('should adapt from JSON',() => {
		let jsonData = [
			{
				"Host":"atreides",
				"Name":"box",
				"Climate":{
					"Temperature":21.99,
					"Humidity":59.46,
					"TemperatureBounds":{
						"Min":20,
						"Max":31
					},
					"HumidityBounds":{
						"Min":40,
						"Max":80
					},
					"ActiveWarnings":1,
					"ActiveEmergencies":1,
					"NumAux":0,
					"Current":{
						"Name":"night",
						"Temperature":22,
						"Humidity":60,
						"Wind":100,
						"VisibleLight":0,
						"UVLight":0},
					"CurrentEnd":null,
					"Next":{
						"Name":"night to day",
						"Temperature":22,
						"Humidity":60,
						"Wind":100,
						"VisibleLight":0,
						"UVLight":0
					},
					"NextEnd":{
						"Name":"day",
						"Temperature":26,
						"Humidity":60,
						"Wind":100,
						"VisibleLight":40,
						"UVLight":100},
					"NextTime":"2021-02-02T06:00:00Z"
				},
				"Stream":null,
				"Alarms":null
			},
			{
				"Host":"fremens",
				"Name":"tunnel",
				"Climate":{
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
					"ActiveEmergencies":1,
					"NumAux":0,
					"Current":{
						"Name":"always-on",
						"Temperature":-1000,
						"Humidity":-1000,
						"Wind":-1000,
						"VisibleLight":100,
						"UVLight":0
					},
					"CurrentEnd":null,
					"Next":null,
					"NextEnd":null,
					"NextTime":null
				},
				"Stream":null,
				"Alarms":[
					{"Reason":"Cannot reach desired humidity","Level":1,"On":true,"Time":"2021-02-01T14:00:04.920610304+01:00"},
					{"Reason":"Cannot reach desired humidity","Level":1,"On":false,"Time":"2021-02-01T14:00:05.420610304+01:00"},
					{"Reason":"Celaeno is empty","Level":2,"On":true,"Time":"2021-02-01T14:02:04.920610304+01:00"},
					{"Reason":"Celaeno is empty","Level":2,"On":false,"Time":"2021-02-01T14:44:06.920610304+01:00"},
					{"Reason":"Celaeno is empty","Level":2,"On":true,"Time":"2021-02-01T14:51:08.920610304+01:00"},
				]
			},
		]
		expect(ZoneReport.adapt(jsonData)).toBeTruthy();

	});

});
