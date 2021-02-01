import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { ZoneSummaryReport } from '@models/zone-summary-report';

import { ZonePreviewComponent } from './zone-preview.component';

describe('ZonePreviewComponent', () => {
	let component: ZonePreviewComponent;
	let fixture: ComponentFixture<ZonePreviewComponent>;
	let summary = new ZoneSummaryReport('somehost','somezone');

	beforeEach(waitForAsync(() => {
		TestBed.configureTestingModule({
			declarations: [ZonePreviewComponent]
		})
			.compileComponents();
	}));



	beforeEach(() => {
		fixture = TestBed.createComponent(ZonePreviewComponent);
		component = fixture.componentInstance;
		component.summary = summary;
		fixture.detectChanges();
	});

	it('should create', () => {
		expect(component).toBeTruthy();
	});

	it('should display zone name',() => {
		const compiled = fixture.debugElement.nativeElement;
		expect(compiled.querySelector('h3').textContent).toContain('somehost.somezone');

	});

});
