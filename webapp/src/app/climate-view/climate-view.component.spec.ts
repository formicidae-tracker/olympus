import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ClimateViewComponent } from './climate-view.component';

import { OlympusService } from '@services/olympus';
import { FakeOlympusService } from '@services/fake-olympus';

describe('ClimateViewComponent', () => {
	let component: ClimateViewComponent;
	let fixture: ComponentFixture<ClimateViewComponent>;
	let olympusFake: any;
	let olympus: any;

	beforeEach(waitForAsync(() => {
		TestBed.configureTestingModule({
			declarations: [ ClimateViewComponent ],
			providers: [{provide: OlympusService, useClass: FakeOlympusService}],
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
