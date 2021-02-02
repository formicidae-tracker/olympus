import { AlarmEvent, AlarmLevel, AlarmReport } from "./alarm";

describe('AlarmEvent', () => {
	describe('empty event',() => {
		let e =  new AlarmEvent();
		it('should create an instance',() => {
			expect(e).toBeTruthy();
		});

		it ('should be off',() => {
			expect(e.on).toBeFalse();
		});

		it ('should have a time',() => {
			expect(e.time).toBeTruthy();
		});
	});

	it('should adapt from JSON',() => {
		let e = AlarmEvent.adapt({
			"Time":"2021-02-04T06:00:00Z",
			"On": true,
		});

		expect(e).toBeTruthy();
		expect(e.time).toEqual(new Date('2021-02-04T06:00:00Z'))
		expect(e.on).toEqual(true);
	});
});

describe('AlarmReport',() => {
	let r: AlarmReport;
	describe('default report',() => {
		beforeEach(() => {
			r = new AlarmReport();
		});

		it('should create an instance',() => {
			expect(r).toBeTruthy();
		});

		it('should have no reason',() => {
			expect(r.reason).toBe('');
		});

		it('should not have a last event',() => {
			expect(r.lastEvent()).toBeFalsy();
		});

		it('should be off',() => {
			expect(r.on()).toBe(false);
		});

		it('should have a last time',() => {
			expect(r.lastTime()).toBeTruthy();
		});

		it('should have an action',() => {
			expect(r.action()).toBe('info');
		});
	});

	describe('off state', () => {
		beforeEach(() => {
			r = new AlarmReport('',AlarmLevel.Warning,[new AlarmEvent(false)]);
		});

		it('should always be info', () => {
			r.level = AlarmLevel.Warning;
			expect(r.action()).toBe('info');
			r.level = AlarmLevel.Critical;
			expect(r.action()).toBe('info');
		});
	});

	describe('on state', () => {
		beforeEach(() => {
			r = new AlarmReport('',AlarmLevel.Warning,[new AlarmEvent(true)]);
		});
		it('should match alarm level',() => {
			let testdata = [
				{ level:AlarmLevel.Warning,action: 'warning'},
				{ level: AlarmLevel.Critical,action: 'danger'},
			];
			for ( let d of testdata ) {
				r.level = d.level;
				expect(r.action()).toEqual(d.action);
			}
		});
	});

	it('should format time since event',() => {
		let start = new Date().getTime();
		let testdata = [
			{time: start - 200,expected: 'now'},
			{time: start - 1200,expected: '1s'},
			{time: start - 60000,expected: '1m'},
			{time: start - 3600000,expected: '1h0m'},
			{time: start - 3600001,expected: '1h0m'},
			{time: start - 3660000,expected: '1h1m'},
		];
		for ( let d of testdata ) {
			let r = new AlarmReport('',AlarmLevel.Warning,[new AlarmEvent(false,new Date(d.time))]);
			expect(r.since(new Date(start))).toEqual(d.expected);
		}
	});

	it('should compare alarm',() => {
		let start = new Date();
		let testdata = [
			{
				a: new AlarmReport(),
				b: new AlarmReport(),
				expected: 0,
			},
			{
				a: new AlarmReport('foo'),
				b: new AlarmReport('bar'),
				expected: 1,
			},
			{
				a: new AlarmReport('foo',AlarmLevel.Critical,[new AlarmEvent(true,new Date(start.getTime()+1))]),
				b: new AlarmReport('bar',AlarmLevel.Critical,[new AlarmEvent(true,new Date(start.getTime()-1))]),
				expected: -1,
			},
			{
				a: new AlarmReport('foo',AlarmLevel.Critical),
				b: new AlarmReport('bar',AlarmLevel.Warning),
				expected: -1,
			},
			{
				a: new AlarmReport('foo',AlarmLevel.Critical,[new AlarmEvent(true,new Date(start.getTime()+1))]),
				b: new AlarmReport('bar',AlarmLevel.Critical),
				expected: -1,
			},

		];
		for ( let d of testdata ) {
			expect(AlarmReport.compare(d.a,d.b)).toBe(d.expected);
		}
	});

});
