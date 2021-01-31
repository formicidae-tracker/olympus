import { Bounds} from './bounds';

export class ZoneClimateStatus {
	constructor(public Temperature: number,
				public Humidity: number,
				public TemperatureBounds: Bounds = null,
				public HumidityBounds: Bounds = null,
				public ActiveWarnings = 0,
				public ActiveEmergencies = 0) {
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
