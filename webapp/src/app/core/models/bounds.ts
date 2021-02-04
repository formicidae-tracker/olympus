export class Bounds {
	private boundary: number;
	constructor(public Min: number = NaN,
				public Max: number = NaN) {
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
		let min: number = NaN;
		let max: number = NaN;
		if (item == null) {
			return new Bounds(min,max);
		}
		if (item.Min != null) {
			min = item.Min;
		}
		if (item.Max != null) {
			max = item.Max;
		}
		return new Bounds(min,max);
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
