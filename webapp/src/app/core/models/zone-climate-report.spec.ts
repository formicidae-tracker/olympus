import { Bounds } from './bounds';
import { ZoneClimateReport} from './zone-climate-report';



describe('ZoneClimateReport', () => {
	it('should create an instance',() =>{
		expect(new ZoneClimateReport(null,0,null,null,null,null,null)).toBeTruthy();
	});

	it('should adapt from null', () => {
		expect(ZoneClimateReport.adapt(null)).toBeNull();
	});

	it('should adapt from full JSON', () => {
		let s = ZoneClimateReport.adapt({
			"Temperature":22.01,
			"Humidity":59.81,
			"TemperatureBounds":{
				"Min":20,
				"Max":31
			},
			"HumidityBounds":{
				"Min":40,
				"Max":80
			},
			"ActiveWarnings":0,
			"ActiveEmergencies":1,
			"NumAux":0,
			"Current":{"Name":"night","Temperature":22,"Humidity":60,"Wind":100,"VisibleLight":0,"UVLight":0},
			"CurrentEnd":null,
			"Next":{"Name":"night to day","Temperature":22,"Humidity":60,"Wind":100,"VisibleLight":0,"UVLight":0},
			"NextEnd":{"Name":"day","Temperature":26,"Humidity":60,"Wind":100,"VisibleLight":40,"UVLight":100},
			"NextTime":"2021-02-04T06:00:00Z"
		});
		expect(s).toBeTruthy();
		expect(s.temperature).toBe(22.01);
		expect(s.temperatureBounds).toEqual(new Bounds(20,31))
		expect(s.humidity).toBe(59.81);
		expect(s.humidityBounds).toEqual(new Bounds(40,80))
		expect(s.numAux).toBe(0);
		expect(s.activeWarnings).toBe(0);
		expect(s.activeEmergencies).toBe(1);
		expect(s.current).toBeTruthy;
		expect(s.currentEnd).toBeNull();
		expect(s.next).toBeTruthy();
		expect(s.nextEnd).toBeTruthy();
		expect(s.nextTime).toBeTruthy();
	});

	it('should adapt from partial JSON', () => {
		let s = ZoneClimateReport.adapt({
			"Temperature":-1000,
			"Humidity":-1000,
			"TemperatureBounds":{
				"Min":null,
				"Max":null
			},"HumidityBounds":{
				"Min":null,
				"Max":null
			},
			"ActiveWarnings":0,
			"ActiveEmergencies":1,
			"NumAux":0,
			"Current":{"Name":"always-on","Temperature":-1000,"Humidity":-1000,"Wind":-1000,"VisibleLight":100,"UVLight":0},
			"CurrentEnd":null,
			"Next":null,
			"NextEnd":null,
			"NextTime":null});
		expect(s.temperature).toBeNaN();
		expect(s.humidity).toBeNaN();
		expect(s.temperatureBounds).toEqual(new Bounds(NaN,NaN));
		expect(s.humidityBounds).toEqual(new Bounds(NaN,NaN));
		expect(s.numAux).toBe(0);
		expect(s.current).toBeTruthy();
		expect(s.currentEnd).toBeNull();
		expect(s.next).toBeNull();
		expect(s.nextEnd).toBeNull();
		expect(s.nextTime).toBeNull();
	});


});
