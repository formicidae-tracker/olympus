import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { Alarm } from '@models/alarm';

import { AlarmComponent } from './alarm.component';

describe('AlarmComponent', () => {
	let component: AlarmComponent;
	let fixture: ComponentFixture<AlarmComponent>;

	beforeEach(waitForAsync(() => {
		TestBed.configureTestingModule({
			declarations: [AlarmComponent]
		})
			.compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(AlarmComponent);
		component = fixture.componentInstance;
		component.alarm = new Alarm('foo',true,new Date(),1,1);
		fixture.detectChanges();
	});

	it('should create', () => {
		expect(component).toBeTruthy();
	});
});
