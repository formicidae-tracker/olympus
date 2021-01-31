import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { ActivatedRoute } from '@angular/router';

import { ZoneComponent } from './zone.component';

import { OlympusService,MockOlympusService } from '@services/olympus';


class MockActivatedRoute {
	static current = {
		hostName: "",
		zoneName: "",
	}
	snapshot = {
		paramMap: {
			get: (key: string) => { return MockActivatedRoute.current[key]; }
		},
	}
}

describe('ZoneComponent', () => {
	let component: ZoneComponent;
	let fixture: ComponentFixture<ZoneComponent>;

	beforeEach(waitForAsync(() => {
		TestBed.configureTestingModule({
			declarations: [ ZoneComponent ],
			providers: [
				{
					provide: OlympusService,
					useClass: MockOlympusService
				},
				{
					provide: ActivatedRoute,
					useClass: MockActivatedRoute,
				},
			],
		})
			.compileComponents();
	}));




	describe('box zone with tracking', () => {
		beforeEach(() => {
			MockActivatedRoute.current = {
				hostName: "somehost",
				zoneName: "box",
			};
			fixture = TestBed.createComponent(ZoneComponent);
			component = fixture.componentInstance;
			fixture.detectChanges();
		});

		it('should create', () => {
			expect(component).toBeTruthy();
		});

	});

	describe('box zone without tracking', () => {
		beforeEach(() => {
			MockActivatedRoute.current = {
				hostName: "notracking",
				zoneName: "box",
			};
			fixture = TestBed.createComponent(ZoneComponent);
			component = fixture.componentInstance;
			fixture.detectChanges();
		});

		it('should create', () => {
			expect(component).toBeTruthy();
		});

	});

	describe('tunnel zone without tracking', () => {
		beforeEach(() => {
			MockActivatedRoute.current = {
				hostName: "somehost",
				zoneName: "tunnel",
			};
			fixture = TestBed.createComponent(ZoneComponent);
			component = fixture.componentInstance;
			fixture.detectChanges();
		});

		it('should create', () => {
			expect(component).toBeTruthy();
		});

	});

	describe('box zone without climate', () => {
		beforeEach(() => {
			MockActivatedRoute.current = {
				hostName: "tracking",
				zoneName: "box",
			};
			fixture = TestBed.createComponent(ZoneComponent);
			component = fixture.componentInstance;
			fixture.detectChanges();
		});

		it('should create', () => {
			expect(component).toBeTruthy();
		});

	});


	describe('unexisting zone', () => {
		beforeEach(() => {
			MockActivatedRoute.current = {
				hostName: "foo",
				zoneName: "bar",
			};
			fixture = TestBed.createComponent(ZoneComponent);
			component = fixture.componentInstance;
			fixture.detectChanges();
		});

		it('should create', () => {
			expect(component).toBeTruthy();
		});

	});


});
