import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { ZoneSummaryReport } from '@models/zone-summary-report';

import { ZonePreviewComponent } from './zone-preview.component';

describe('ZonePreviewComponent', () => {
	let component: ZonePreviewComponent;
	let fixture: ComponentFixture<ZonePreviewComponent>;

	beforeEach(waitForAsync(() => {
		TestBed.configureTestingModule({
			declarations: [ZonePreviewComponent]
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
			expect(res.textContent).toBe('No current state');

		});

	});



});
