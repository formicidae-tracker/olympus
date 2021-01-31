import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { convertToParamMap,ActivatedRoute } from '@angular/router';
import { of } from 'rxjs';

import { ZoneComponent } from './zone.component';

import { OlympusService,MockOlympusService } from '@services/olympus';


class MockActivatedRoute {
	public paramMap = of(convertToParamMap({
	}));
}

describe('ZoneComponent', () => {
	let component: ZoneComponent;
	let fixture: ComponentFixture<ZoneComponent>;

	let activatedRouteMock = {
		snapshot: {
			paramMap:convertToParamMap({
				hostName: "somehost",
				zoneName: "box",
			}),
		},
	};

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
					useValue: activatedRouteMock,
				},
			],
		})
			.compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(ZoneComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it('should create', () => {
		expect(component).toBeTruthy();
	});
});
