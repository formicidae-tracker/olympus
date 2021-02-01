
export class Bounds {
	private boundary: number;
	constructor(public Min: number,
				public Max: number) {
		this.boundary = 0;
		if (isNaN(this.Min) == false
			&& isNaN(this.Max) == false ) {
			if ( this.Min > this.Max ) {
				let tmp = this.Min;
				this.Min = this.Max;
				this.Max = tmp;
			}
			this.boundary = 0.05 * (this.Max-this.Min);
		}
	}

	static adapt(item: any): Bounds {
		let res = new Bounds(NaN,NaN);
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

	status(v: number): string {
		if ( isNaN(this.Min) == false ) {
			if ( v < this.Min ) {
				return 'danger';
			}
			if ( v < this.Min + this.boundary ) {
				return 'warning';
			}
		}

		if ( isNaN(this.Max) == false ) {
			if ( v > this.Max ) {
				return 'danger';
			}
			if ( v > this.Max - this.boundary ) {
				return 'warning';
			}
		}
		return 'success';
	}
}
