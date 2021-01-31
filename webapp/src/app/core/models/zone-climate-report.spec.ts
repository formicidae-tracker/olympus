import { ZoneClimateReport} from './zone-climate-report';



describe('ZoneClimateReport', () => {
	it('should create an instance',() =>{
		expect(new ZoneClimateReport(null,0,null,null,null,null,null)).toBeTruthy();
	});

	it('should adapt from JSON', () => {
		let s = ZoneClimateReport.adapt({"Temperature":20.0,
										 "Humidity":12.1,
										 "TemperatureBounds":{
											 "Min":null,
											 "Max":null
										 },
										 "HumidityBounds":{
											 "Min":null,
											 "Max":null},
										 "ActiveWarnings":0,
										 "ActiveEmergencies":0,
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
										})
		expect(s).toBeTruthy();
		expect(s.ClimateStatus.Temperature).toBe(20);
		expect(s.ClimateStatus.Humidity).toBe(12.1);
		expect(s.CurrentEnd).toBe(null);
		expect(s.Next).toBe(null);
		expect(s.NextEnd).toBe(null);
		expect(s.NextTime).toBe(null);
	});


});
