export class State {
	public Name: string;
	public Humidity: number;
	public Temperature: number;
	public Wind: number;
	public VisibleLight: number;
	public UVLight: number;
	constructor(name: string,
				humidity: number,
				temperature:number,
				wind: number,
				visibleLight: number,
				uvLight: number) {
		this.Name = name;
		this.Humidity = this.checkUndefined(humidity);
		this.Temperature = this.checkUndefined(temperature);
		this.Wind = this.checkUndefined(wind);
		this.VisibleLight = this.checkUndefined(visibleLight);
		this.UVLight = this.checkUndefined(uvLight);
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
