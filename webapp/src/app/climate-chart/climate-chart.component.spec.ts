import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ClimateChartComponent } from './climate-chart.component';

import { OlympusService,MockOlympusService } from '@services/olympus';

describe('ClimateChartComponent', () => {
	let component: ClimateChartComponent;
	let fixture: ComponentFixture<ClimateChartComponent>;

	beforeEach(waitForAsync(() => {
		TestBed.configureTestingModule({
			declarations: [ ClimateChartComponent ],
			providers: [{provide: OlympusService, useClass: MockOlympusService}],
		})
			.compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(ClimateChartComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it('should create', () => {
		expect(component).toBeTruthy();
	});
});
