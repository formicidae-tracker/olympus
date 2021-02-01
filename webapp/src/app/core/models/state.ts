export class State {
	public name: string;
	public humidity: number;
	public temperature: number;
	public wind: number;
	public visibleLight: number;
	public uvLight: number;
	constructor(name: string,
				humidity: number,
				temperature:number,
				wind: number,
				visibleLight: number,
				uvLight: number) {
		this.name = name;
		this.humidity = this.checkUndefined(humidity);
		this.temperature = this.checkUndefined(temperature);
		this.wind = this.checkUndefined(wind);
		this.visibleLight = this.checkUndefined(visibleLight);
		this.uvLight = this.checkUndefined(uvLight);
	}

	private checkUndefined(v: number): number{
		if ( v <= -1000.0 ) {
			return NaN;
		}
		return v;
	}

	static adapt(item: any): State {
		if ( item == null ) {
			return null;
		}
		return new State(item.Name,
						 item.Humidity,
						 item.Temperature,
						 item.Wind,
						 item.VisibleLight,
						 item.UVLight);
	}

}
