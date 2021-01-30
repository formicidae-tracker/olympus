import { Bounds,BoundsAdapter } from './bounds';

describe('Bounds', () => {
	it('should create an instance', () => {
		expect(new Bounds(0,100)).toBeTruthy();
	});

	it('should always represent a valid range',() => {
		let b = new Bounds(29,1);
		expect(b.Min < b.Max).toBeTrue();
	});

	let adapter = new BoundsAdapter();
	it('should adapt a default from null', () => {
		let b = adapter.adapt(null);

		expect(b).toBeTruthy();
		expect(b.Min).toBe(0);
		expect(b.Max).toBe(100);
	});


	it('should adapt from partial definition', () => {
		let b = adapter.adapt({"Min": 42});
		expect(b.Min).toBe(42);
		expect(b.Max).toBe(100);

		b = adapter.adapt({"Max": 42});
		expect(b.Min).toBe(0);
		expect(b.Max).toBe(42);
	});

	it('should adapt from full definition',() => {
		let b = adapter.adapt({"Min":10,"Max":20});
		expect(b.Min).toBe(10);
		expect(b.Max).toBe(20);
	});

});
