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

	hasTemperature(): boolean {
		return this.temperature.length > 0 || this.temperatureAux.length > 0;
	}

	hasHumidity(): boolean {
		return this.humidity.length > 0;
	}

	hasData(): boolean {
		return this.hasTemperature() || this.hasHumidity();
	}

}
