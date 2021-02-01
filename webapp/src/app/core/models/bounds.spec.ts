import { Bounds } from './bounds';

describe('Bounds', () => {
	it('should create an instance', () => {
		expect(new Bounds(0,100)).toBeTruthy();
	});

	it('should always represent a valid range',() => {
		let b = new Bounds(29,1);
		expect(b.Min < b.Max).toBeTrue();
	});

	it('should adapt a default from null', () => {
		let b = Bounds.adapt(null);

		expect(b).toEqual(new Bounds(NaN,NaN));
		expect(b.Max).toBeNaN();
	});


	it('should adapt from partial definition', () => {
		let b = Bounds.adapt({"Min": 42});
		expect(b).toEqual(new Bounds(42,NaN));

		b = Bounds.adapt({"Max": 42});
		expect(b).toEqual(new Bounds(NaN,42));

	});

	it('should adapt from full definition',() => {
		let b = Bounds.adapt({"Min":10,"Max":20});
		expect(b).toEqual(new Bounds(10,20));
	});

});
