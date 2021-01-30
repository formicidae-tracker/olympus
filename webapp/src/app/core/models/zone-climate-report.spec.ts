import { ZoneClimateReport,ZoneClimateReportAdapter} from './zone-climate-report';

import { StateAdapter } from './state';
import { ZoneClimateStatusAdapter } from './zone-climate-status';
import { BoundsAdapter } from './bounds';


describe('ZoneClimateReport', () => {
	let adapter: ZoneClimateReportAdapter;

	beforeEach(() => {
		adapter = new ZoneClimateReportAdapter(new ZoneClimateStatusAdapter(new BoundsAdapter()),
											   new StateAdapter());
	});

	it('should create an instance',() =>{
		expect(new ZoneClimateReport(null,0,null,null,null,null,null)).toBeTruthy();
	});

	it('should adapt from JSON', () => {
		let s = adapter.adapt({"Temperature":20.0,
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
		expect(s.Temperature).toBe(20);
		expect(s.Humidity).toBe(12.1);
		expect(s.CurrentEnd).toBe(null);
		expect(s.Next).toBe(null);
		expect(s.NextEnd).toBe(null);
		expect(s.NextTime).toBe(null);
	});


});
