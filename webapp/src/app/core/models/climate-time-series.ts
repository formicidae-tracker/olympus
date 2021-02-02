export class ClimateTimeSeries {
	constructor(public humidity: any[] = [],
				public temperature: any[] = [],
				public temperatureAux: any[][] =[]){}

	static adapt(item: any): ClimateTimeSeries {
		let humidity: any[] = [];
		let temperature: any[] = [];
		let tempAux: any[][] = [];
		if ( item.Humidity != null ) {
			humidity = item.Humidity;
		}
		if ( item.TemperatureAnt != null ) {
			temperature = item.TemperatureAnt;
		}
		if ( item.TemperatureAux != null ) {
			tempAux = item.TemperatureAux;
		}

		return new ClimateTimeSeries(humidity,temperature,tempAux);
	}
}
