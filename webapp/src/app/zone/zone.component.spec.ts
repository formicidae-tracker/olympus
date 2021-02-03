import { ComponentFixture, TestBed, waitForAsync,fakeAsync,tick,discardPeriodicTasks } from '@angular/core/testing';
import { ActivatedRoute, convertToParamMap } from '@angular/router';

import { ZoneComponent } from './zone.component';
import { of } from 'rxjs';
import { OlympusService } from '@services/olympus';
import { FakeOlympusService } from '@services/fake-olympus';



describe('ZoneComponent', () => {
	let component: ZoneComponent;
	let fixture: ComponentFixture<ZoneComponent>;
	let fakeOlympus = new FakeOlympusService();
	let olympus: any;
	let zoneReportSpy: any;
	let route = {
		paramMap: of(convertToParamMap({
			hostName: '',
			zoneName: '',
		})),
	}

	beforeEach(waitForAsync(() => {
		olympus = jasmine.createSpyObj('OlympusService',['zoneReport']);
		zoneReportSpy = olympus.zoneReport.and.callFake(function (host,name) {
			return fakeOlympus.zoneReport(host,name);
		});
		TestBed.configureTestingModule({
			declarations: [ ZoneComponent ],
			providers: [
				{provide: OlympusService, useValue: olympus},
				{provide: ActivatedRoute, useValue: route},
			],
		}).compileComponents();
	}));


	describe('box zone with tracking', () => {
		beforeEach(fakeAsync(() => {
			route.paramMap = of(convertToParamMap({
				hostName: 'somehost',
				zoneName: 'box',
			}));

			fixture = TestBed.createComponent(ZoneComponent);
			component = fixture.componentInstance;
			fixture.detectChanges();
			tick();
			fixture.detectChanges();
			discardPeriodicTasks();
		}));


		it('should create with the right parameters',() => {

			expect(zoneReportSpy.calls.any()).toBe(true,'zoneReport called');

			expect(component).toBeTruthy();
			expect(component.hostName).toBe('somehost','private hostName');
			expect(component.zoneName).toBe('box','private zoneName');
			expect(component.zone.host).toBe('somehost','zone.host');
			expect(component.zone.name).toBe('box','zone.name');
		});

		it('should display the alarm logs', () => {
			const compiled = fixture.debugElement.nativeElement;
			const selection = compiled.querySelectorAll('.sidebar-sticky app-alarm-list');
			expect(selection.length).toBe(1,compiled);
		});

		it('should display the video stream',() => {
			const compiled = fixture.debugElement.nativeElement;
			const selection = compiled.querySelectorAll('app-video-js');
			expect(selection.length).toBe(1);
		});

		it('should display the climate state',() => {
			const compiled = fixture.debugElement.nativeElement;
			let selection = compiled.querySelectorAll('main h2');
			expect(selection.length).toBe(1)
			if ( selection.length == 0 ) {
				return;
			}
			expect(selection[0].textContent).toContain('Climate States for somehost.box');
			selection = compiled.querySelectorAll('app-climate-view');
			expect(selection.length).toBe(1);
		});

		it('should remove the loading screen once finished',()=> {

			expect(component.loading).toBe(false,'loading is finished');
			expect(component.unavailable()).toBe(false,'is unavailable');

			const compiled = fixture.debugElement.nativeElement;
			const selection = compiled.querySelectorAll('main .jumbotron h1');
			expect(selection.length).toBe(0);
		});


	});

	describe('box zone without tracking', () => {
		beforeEach(fakeAsync(() => {
			route.paramMap = of(convertToParamMap({
				hostName: 'notracking',
				zoneName: 'box',
			}));
			fixture = TestBed.createComponent(ZoneComponent);
			component = fixture.componentInstance;
			fixture.detectChanges();
			tick();
			fixture.detectChanges();
			discardPeriodicTasks();
		}));

		it('should create with the right parameters', () => {
			expect(component).toBeTruthy();
			expect(component.hostName).toBe('notracking','private hostName');
			expect(component.zoneName).toBe('box','private zoneName');
			expect(component.zone.host).toBe('notracking','zone.host');
			expect(component.zone.name).toBe('box','zone.name');
		});

		it('should display the alarm logs', () => {
			const compiled = fixture.debugElement.nativeElement;
			const selection = compiled.querySelectorAll('.sidebar-sticky app-alarm-list');
			expect(selection.length).toBe(1);
		});

		it('should display the climate state',() => {
			const compiled = fixture.debugElement.nativeElement;
			let selection = compiled.querySelectorAll('main h2');
			expect(selection.length).toBe(1)
			if ( selection.length == 0 ) {
				return;
			}
			expect(selection[0].textContent).toContain('Climate States for notracking.box');
			selection = compiled.querySelectorAll('app-climate-view');
			expect(selection.length).toBe(1);
		});

		it('should not display the video stream',() => {
			const compiled = fixture.debugElement.nativeElement;
			const selection = compiled.querySelectorAll('app-video-js');
			expect(selection.length).toBe(0);
		});

		it('should remove the loading screen once finished',()=> {

			expect(component.loading).toBe(false,'loading is finished');
			expect(component.unavailable()).toBe(false,'is unavailable');

			const compiled = fixture.debugElement.nativeElement;
			const selection = compiled.querySelectorAll('main .jumbotron h1');
			expect(selection.length).toBe(0);
		});

	});

	describe('tunnel zone without tracking', () => {
		beforeEach(fakeAsync(() => {
			route.paramMap = of(convertToParamMap({
				hostName: 'somehost',
				zoneName: 'tunnel',
			}));
			fixture = TestBed.createComponent(ZoneComponent);
			component = fixture.componentInstance;
			fixture.detectChanges();
			tick();
			fixture.detectChanges();
			discardPeriodicTasks();
		}));

		it('should create with the right parameters', () => {
			expect(component).toBeTruthy();
			expect(component.hostName).toBe('somehost','private hostName');
			expect(component.zoneName).toBe('tunnel','private zoneName');
			expect(component.zone.host).toBe('somehost','zone.host');
			expect(component.zone.name).toBe('tunnel','zone.name');
		});

		it('should display the alarm logs', () => {
			const compiled = fixture.debugElement.nativeElement;
			const selection = compiled.querySelectorAll('.sidebar-sticky app-alarm-list');
			expect(selection.length).toBe(1);
		});

		it('should display the climate state',() => {
			const compiled = fixture.debugElement.nativeElement;
			let selection = compiled.querySelectorAll('main h2');
			expect(selection.length).toBe(1)
			if ( selection.length == 0 ) {
				return;
			}
			expect(selection[0].textContent).toContain('Climate States for somehost.tunnel');
			selection = compiled.querySelectorAll('app-climate-view');
			expect(selection.length).toBe(1);
		});

		it('should not display the video stream',() => {
			const compiled = fixture.debugElement.nativeElement;
			const selection = compiled.querySelectorAll('app-video-js');
			expect(selection.length).toBe(0);
		});

		it('should remove the loading screen once finished',()=> {
			expect(component.loading).toBe(false,'loading is finished');
			expect(component.unavailable()).toBe(false,'is unavailable');

			const compiled = fixture.debugElement.nativeElement;
			const selection = compiled.querySelectorAll('main .jumbotron h1');
			expect(selection.length).toBe(0);
		});

	});

	describe('box zone without climate but tracking', () => {
		beforeEach(fakeAsync(() => {
			route.paramMap = of(convertToParamMap({
				hostName: 'onlytracking',
				zoneName: 'box',
			}));
			fixture = TestBed.createComponent(ZoneComponent);
			component = fixture.componentInstance;
			fixture.detectChanges();
			tick();
			fixture.detectChanges();
			discardPeriodicTasks();
		}));

		it('should create with the right parameters', () => {
			expect(component).toBeTruthy();
			expect(component.hostName).toBe('onlytracking','private hostName');
			expect(component.zoneName).toBe('box','private zoneName');
			expect(component.zone.host).toBe('onlytracking','zone.host');
			expect(component.zone.name).toBe('box','zone.name');
		});

		it('should not display the alarm logs', () => {
			const compiled = fixture.debugElement.nativeElement;
			const selection = compiled.querySelectorAll('.sidebar-sticky app-alarm-list');
			expect(selection.length).toBe(0,compiled.innerHTML);
		});

		it('should not display the climate state',() => {
			const compiled = fixture.debugElement.nativeElement;
			let selection = compiled.querySelectorAll('main h2');
			expect(selection.length).toBe(0)
			selection = compiled.querySelectorAll('app-climate-view');
			expect(selection.length).toBe(0);
		});


		it('should display the video stream',() => {
			const compiled = fixture.debugElement.nativeElement;
			const selection = compiled.querySelectorAll('app-video-js');
			expect(selection.length).toBe(1);
		});

		it('should remove the loading screen once finished',()=> {

			expect(component.loading).toBe(false,'loading is finished');
			expect(component.unavailable()).toBe(false,'is unavailable');

			const compiled = fixture.debugElement.nativeElement;
			const selection = compiled.querySelectorAll('main .jumbotron h1');
			expect(selection.length).toBe(0);
		});

	});


	describe('unexisting zone', () => {
		beforeEach(fakeAsync(() => {
			route.paramMap = of(convertToParamMap({
				hostName: 'foo',
				zoneName: 'bar',
			}));
			fixture = TestBed.createComponent(ZoneComponent);
			component = fixture.componentInstance;
			fixture.detectChanges();
			tick();
			fixture.detectChanges();
			discardPeriodicTasks();
		}));

		it('should create with the right parameters', () => {
			expect(component).toBeTruthy();
			expect(component.hostName).toBe('foo','private hostName');
			expect(component.zoneName).toBe('bar','private zoneName');
		});

		it('should have no zone object',() => {
			expect(component.unavailable()).toBeTrue();
			expect(component.zone).toBeFalsy();
		});

		it('should display zone unavailble',() => {
			const compiled = fixture.debugElement.nativeElement;
			const selection = compiled.querySelectorAll('main .jumbotron h1');
			expect(selection.length).toBe(1);
			expect(selection[0].textContent).toContain('foo.bar is unavailable');
		});

		it('should not display the alarm logs', () => {
			const compiled = fixture.debugElement.nativeElement;
			const selection = compiled.querySelectorAll('.sidebar-sticky app-alarm-list');
			expect(selection.length).toBe(0,compiled.innerHTML);
		});

		it('should not display the climate state',() => {
			const compiled = fixture.debugElement.nativeElement;
			let selection = compiled.querySelectorAll('main h2');
			expect(selection.length).toBe(0)
			selection = compiled.querySelectorAll('app-climate-view');
			expect(selection.length).toBe(0);
		});

		it('should not display the video stream',() => {
			const compiled = fixture.debugElement.nativeElement;
			const selection = compiled.querySelectorAll('app-video-js');
			expect(selection.length).toBe(0);
		});

		it('should remove the loading screen once finished',()=> {
			expect(component.loading).toBe(false,'loading is finished');
			expect(component.unavailable()).toBe(true,'is unavailable');

			const compiled = fixture.debugElement.nativeElement;
			const selection = compiled.querySelectorAll('main .jumbotron h1');
			expect(selection.length).toBe(1);

			expect(selection[0].textContent).not.toContain('Loading foo.bar');
		});

	});

	describe('loading screen',() => {
		beforeEach(fakeAsync(() => {
			route.paramMap = of(convertToParamMap({
				hostName: 'foo',
				zoneName: 'bar',
			}));
			fixture = TestBed.createComponent(ZoneComponent);
			component = fixture.componentInstance;
			fixture.detectChanges();
			discardPeriodicTasks();
		}));

		it('should display loading until service subscription finishes',fakeAsync(()=> {

			const compiled = fixture.debugElement.nativeElement;
			const selection = compiled.querySelectorAll('main .jumbotron h1');
			expect(selection.length).toBe(1);

			expect(selection[0].textContent).toContain('Loading foo.bar');
			expect(compiled.querySelector('main .jumbotron p').textContent).toContain('Firefox is known to be quite slow to load this page.');


		}));

		it('should not display the video stream',() => {
			const compiled = fixture.debugElement.nativeElement;
			const selection = compiled.querySelectorAll('app-video-js');
			expect(selection.length).toBe(0);
		});

		it('should not display the alarm logs', () => {
			const compiled = fixture.debugElement.nativeElement;
			const selection = compiled.querySelectorAll('.sidebar-sticky app-alarm-list');
			expect(selection.length).toBe(0,compiled.innerHTML);
		});

		it('should not display the climate state',() => {
			const compiled = fixture.debugElement.nativeElement;
			let selection = compiled.querySelectorAll('main h2');
			expect(selection.length).toBe(0)
			selection = compiled.querySelectorAll('app-climate-view');
			expect(selection.length).toBe(0);
		});

	});


});
