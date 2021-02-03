import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { ActivatedRoute,convertToParamMap } from '@angular/router';

import { ZoneComponent, ZoneState } from './zone.component';
import { of } from 'rxjs';
import { OlympusService,MockOlympusService } from '@services/olympus';



describe('ZoneComponent', () => {
	let component: ZoneComponent;
	let fixture: ComponentFixture<ZoneComponent>;
	let olympus = new MockOlympusService();
	let route = {
		paramMap: of(convertToParamMap({
			hostName: '',
			zoneName: '',
		})),
	}
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
					useValue: route,
				},
			],
		})
			.compileComponents();
	}));


	describe('box zone with tracking', () => {
		beforeEach(() => {
			route.paramMap = of(convertToParamMap({
				hostName: 'somehost',
				zoneName: 'box',
			}));
			fixture = TestBed.createComponent(ZoneComponent);
			component = fixture.componentInstance;
			component.zone = olympus.zoneReportStatic('somehost','box');
			component.state = ZoneState.Loaded;
			fixture.detectChanges();
		});

		it('should create with the right parameters', () => {
			expect(component).toBeTruthy();
			expect(component.hostName).toBe('somehost');
			expect(component.zoneName).toBe('box');
		});

		it('should have the right zone',() => {
			expect(component.zone.host).toBe('somehost');
		})

	});

	describe('box zone without tracking', () => {
		beforeEach(() => {
			route.paramMap = of(convertToParamMap({
				hostName: 'notracking',
				zoneName: 'box',
			}));
			fixture = TestBed.createComponent(ZoneComponent);
			component = fixture.componentInstance;
			fixture.detectChanges();
		});

		it('should create', () => {
			expect(component).toBeTruthy();
			expect(component.hostName).toBe('notracking');
			expect(component.zoneName).toBe('box');
		});

	});

	describe('tunnel zone without tracking', () => {
		beforeEach(() => {
			route.paramMap = of(convertToParamMap({
				hostName: 'somehost',
				zoneName: 'tunnel',
			}));
			fixture = TestBed.createComponent(ZoneComponent);
			component = fixture.componentInstance;
			fixture.detectChanges();
		});

		it('should create', () => {
			expect(component).toBeTruthy();
			expect(component.hostName).toBe('somehost');
			expect(component.zoneName).toBe('tunnel');

		});

	});

	describe('box zone without climate', () => {
		beforeEach(() => {
			route.paramMap = of(convertToParamMap({
				hostName: 'onlytracking',
				zoneName: 'box',
			}));
			fixture = TestBed.createComponent(ZoneComponent);
			component = fixture.componentInstance;
			fixture.detectChanges();
		});

		it('should create', () => {
			expect(component).toBeTruthy();
			expect(component.hostName).toBe('onlytracking');
			expect(component.zoneName).toBe('box');
		});

	});


	describe('unexisting zone', () => {
		beforeEach(() => {
			route.paramMap = of(convertToParamMap({
				hostName: 'foo',
				zoneName: 'bar',
			}));
			fixture = TestBed.createComponent(ZoneComponent);
			component = fixture.componentInstance;
			fixture.detectChanges();
		});

		it('should create', () => {
			expect(component).toBeTruthy();
			expect(component.hostName).toBe('foo');
			expect(component.zoneName).toBe('bar');

		});

	});


});
