import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { NgbCollapse, NgbModule } from '@ng-bootstrap/ng-bootstrap';

import { AlarmComponent } from './alarm.component';

describe('AlarmComponent', () => {
	let component: AlarmComponent;
	let fixture: ComponentFixture<AlarmComponent>;

	beforeEach(waitForAsync(() => {
		TestBed.configureTestingModule({
			imports: [ NgbModule ],
			declarations: [AlarmComponent],
		})
			.compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(AlarmComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it('should create', () => {
		expect(component).toBeTruthy();
	});
});
