import {State} from './state';

describe('State' , () => {

	it('should create an instance', () => {
		expect(new State('foo',12,23,56,0,10)).toBeTruthy();
	});

	it('should sets back undefined', () => {
		let s = new State('foo',-1000,NaN,-1000,-1001,-1000);
		expect(s.Humidity).toBeNaN();
		expect(s.Temperature).toBeNaN();
		expect(s.Wind).toBeNaN();
		expect(s.VisibleLight).toBeNaN();
		expect(s.UVLight).toBeNaN();
	})

	it('should adapt from json', () => {
		let s = State.adapt({"Name":"always-on",
							 "Temperature":-1000,
							 "Humidity":-1000,
							 "Wind":-1000,
							 "VisibleLight":100,
							 "UVLight":0});
		expect(s.Name).toBe('always-on');
		expect(s.Temperature).toBeNaN();
		expect(s.Humidity).toBeNaN();
		expect(s.Wind).toBeNaN();
		expect(s.VisibleLight).toBe(100);
		expect(s.UVLight).toBe(0);
	});

});
