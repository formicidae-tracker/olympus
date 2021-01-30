import { Injectable } from '@angular/core';
import { Adapter } from './adapter';


export class Bounds {
	constructor(public Min: number,
				public Max: number) {}
}

@Injectable({
    providedIn: 'root'
})

export class BoundsAdapter implements Adapter<Bounds> {
	adapt(item: any): Bounds {
		let res = new Bounds(0,100);
		if (item == null) {
			return res;
		}
		if (item.Min != null) {
			res.Min = item.Min;
		}
		if (item.Max != null) {
			res.Max = item.Max;
		}
		return res;
	}
}
