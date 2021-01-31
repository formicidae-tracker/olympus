export class ClimateTimeSeries {
	constructor(public Humidity: any[],
				public Temperature: any[],
				public TemperatureAux: any[][]){}

	static adapt(item: any): ClimateTimeSeries {
		return new ClimateTimeSeries(item.Humidity,
									 item.TemperatureAnt,
									 item.TemperatureAux);
	}
}
