import { Injectable } from '@angular/core';
import { Adapter } from './adapter';

import { Bounds,BoundsAdapter} from './bounds';

export class ZoneClimateStatus {
	constructor(public Temperature: number,
				public Humidity: number,
				public TemperatureBounds: Bounds = null,
				public HumidityBounds: Bounds = null,
				public ActiveWarnings = 0,
				public ActiveEmergencies = 0) {
	}
}

@Injectable({
    providedIn: 'root'
})

export class ZoneClimateStatusAdapter implements Adapter<ZoneClimateStatus> {
	constructor(private boundsAdapter: BoundsAdapter) {}
	adapt(item: any) : ZoneClimateStatus {
		return new ZoneClimateStatus(item.Temperature,
									 item.Humidity,
									 this.boundsAdapter.adapt(item.TemperatureBounds),
									 this.boundsAdapter.adapt(item.HumidityBounds),
									 item.ActiveWarnings,
									 item.ActiveEmergencies);
	}
}
