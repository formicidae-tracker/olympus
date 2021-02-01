import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { Bounds } from '@models/bounds';
import { StreamInfo } from '@models/stream-info';
import { ZoneClimateReport } from '@models/zone-climate-report';
import { ZoneSummaryReport } from '@models/zone-summary-report';
import { MockOlympusService } from '@services/olympus';

import { ZonePreviewComponent } from './zone-preview.component';

describe('ZonePreviewComponent', () => {
	let component: ZonePreviewComponent;
	let fixture: ComponentFixture<ZonePreviewComponent>;
	let olympus: MockOlympusService = new MockOlympusService();
	beforeEach(waitForAsync(() => {
		TestBed.configureTestingModule({
			imports: [
				RouterTestingModule,
			],
			declarations: [
				ZonePreviewComponent,
			]
		})
			.compileComponents();
	}));


	beforeEach(() => {
		fixture = TestBed.createComponent(ZonePreviewComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it('should create', () => {
		expect(component).toBeTruthy();
	});


	describe('Hanging zone',() => {

		beforeEach( () => {
			component.summary = new ZoneSummaryReport('hanging','box')
			fixture.detectChanges();
		});

		it('should display zone name',() => {
			const compiled = fixture.debugElement.nativeElement;
			expect(compiled.querySelector('h3').textContent).toContain('hanging.box');
		});

		it('should display a tracking icon', () => {
			const compiled = fixture.debugElement.nativeElement;
			expect(compiled.querySelector('svg.card-img-top')).toBeTruthy();
			expect(compiled.querySelector('img.card-img-top')).toBeFalsy();
		});

		it('should display no current state', () => {
			const compiled = fixture.debugElement.nativeElement;
			let res = compiled.querySelector('li')
			expect(res.textContent).toContain('No climate control');

		});

	});


	describe('Tracking only',() => {

		beforeEach( () => {
			component.summary = olympus.zoneSummariesStatic()[1];
			fixture.detectChanges();
		});

		it('should display zone name',() => {
			const compiled = fixture.debugElement.nativeElement;
			expect(compiled.querySelector('h3').textContent).toContain('onlytracking.box');
		});

		it('should display a tracking thumbnail', () => {
			const compiled = fixture.debugElement.nativeElement;
			expect(compiled.querySelector('svg.card-img-top')).toBeFalsy();
			expect(compiled.querySelector('img.card-img-top')
				.attributes.getNamedItem('src').value).toBe('/olympus/onlytracking.png');
		});

		it('should display no current state', () => {
			const compiled = fixture.debugElement.nativeElement;
			let res = compiled.querySelector('li')
			expect(res.textContent).toContain('No climate control');

		});

	});

	describe('Tracking and climate',() => {

		beforeEach( () => {
			component.summary = olympus.zoneSummariesStatic()[2];
			fixture.detectChanges();
		});

		it('should display zone name',() => {
			const compiled = fixture.debugElement.nativeElement;
			expect(compiled.querySelector('h3').textContent).toContain('somehost.box');
		});

		it('should display a tracking icon', () => {
			const compiled = fixture.debugElement.nativeElement;
			expect(compiled.querySelector('svg.card-img-top')).toBeFalsy();
			expect(compiled.querySelector('img.card-img-top')
				.attributes.getNamedItem('src').value).toBe('/olympus/somehost.png');
		});

		it('should display no current state', () => {
			const compiled = fixture.debugElement.nativeElement;
			let res = compiled.querySelector('li')
			expect(res.textContent).toContain("Current state: 'day'Next state: 'day to night' at ");
		});

	});



});
