import { ComponentFixture, TestBed } from '@angular/core/testing';
import { FakeOlympusService } from '@services/fake-olympus';
import { OlympusService } from '@services/olympus';

import { LogsComponent } from './logs.component';

describe('LogsComponent', () => {
	let component: LogsComponent;
	let fixture: ComponentFixture<LogsComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [LogsComponent],
			providers: [
				{ provide: OlympusService, useClass: FakeOlympusService },
				]
		})
			.compileComponents();
	});

	beforeEach(() => {
		fixture = TestBed.createComponent(LogsComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it('should create', () => {
		expect(component).toBeTruthy();
	});
});
