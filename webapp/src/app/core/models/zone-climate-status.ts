import { Bounds} from './bounds';

export class ZoneClimateStatus {
	constructor(public Temperature: number,
				public Humidity: number,
				public TemperatureBounds: Bounds = new Bounds(NaN,NaN),
				public HumidityBounds: Bounds = new Bounds(NaN,NaN),
				public ActiveWarnings = 0,
				public ActiveEmergencies = 0) {
		if ( this.Temperature <= -1000.0 ) {
			this.Temperature = NaN;
		}
		if ( this.Humidity <= -1000.0 ) {
			this.Humidity = NaN
		}
	}
	static adapt(item: any) : ZoneClimateStatus {
		if ( item == null ) {
			return null;
		}

		return new ZoneClimateStatus(item.Temperature,
									 item.Humidity,
									 Bounds.adapt(item.TemperatureBounds),
									 Bounds.adapt(item.HumidityBounds),
									 item.ActiveWarnings,
									 item.ActiveEmergencies);
	}
}
