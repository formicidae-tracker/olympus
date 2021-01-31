
export class Bounds {
	constructor(public Min: number,
				public Max: number) {
		if (this.Min > this.Max ) {
			let tmp = this.Min;
			this.Min = this.Max;
			this.Max = tmp;
		}
	}

	static adapt(item: any): Bounds {
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
