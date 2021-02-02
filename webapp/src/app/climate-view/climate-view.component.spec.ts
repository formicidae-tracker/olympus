import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ClimateViewComponent } from './climate-view.component';

import { OlympusService,MockOlympusService } from '@services/olympus';

describe('ClimateViewComponent', () => {
	let component: ClimateViewComponent;
	let fixture: ComponentFixture<ClimateViewComponent>;

	beforeEach(waitForAsync(() => {
		TestBed.configureTestingModule({
			declarations: [ ClimateViewComponent ],
			providers: [{provide: OlympusService, useClass: MockOlympusService}],
		})
			.compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(ClimateViewComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it('should create', () => {
		expect(component).toBeTruthy();
	});
});
