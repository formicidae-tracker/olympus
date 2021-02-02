import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { ActivatedRoute } from '@angular/router';

import { ZoneComponent } from './zone.component';
import { of } from 'rxjs';
import { OlympusService,MockOlympusService } from '@services/olympus';



describe('ZoneComponent', () => {
	let component: ZoneComponent;
	let fixture: ComponentFixture<ZoneComponent>;
	let olympus = new MockOlympusService();
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
					useValue: {

					},
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
			component.zone = olympus.zoneReportStatic('somehost','box');
			fixture.detectChanges();
		});

		it('should create', () => {
			expect(component).toBeTruthy();
		});

		it('should have the right zone',() => {
			expect(component.hostName).toBe('somehost');
			expect(component.zoneName).toBe('box');
		})

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
