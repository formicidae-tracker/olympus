import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { FakeOlympusService } from '@services/fake-olympus';

import { OlympusService } from '@services/olympus';

import { HomeComponent } from './home.component';

describe('HomeComponent', () => {
	let component: HomeComponent;
	let fixture: ComponentFixture<HomeComponent>;

	beforeEach(waitForAsync(() => {
		TestBed.configureTestingModule({
			declarations: [ HomeComponent ],
			providers: [
				{ provide: OlympusService, useClass: FakeOlympusService },
			],

		})
			.compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(HomeComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it('should create', () => {
		expect(component).toBeTruthy();
	});
});
