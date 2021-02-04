import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { State } from '@models/state';

import { StateComponent } from './state.component';

describe('StateComponent', () => {
	let component: StateComponent;
	let fixture: ComponentFixture<StateComponent>;

	beforeEach(waitForAsync(() => {
		TestBed.configureTestingModule({
			declarations: [StateComponent]
		})
			.compileComponents();
	}));


	describe('display value', () => {
		it('renames undefined values',() => {
			expect(StateComponent.displayValue(NaN)).toBe('N.A.');
			expect(StateComponent.displayValue(null)).toBe('N.A.');
			expect(StateComponent.displayValue(0.0)).not.toBe('N.A.');
		});

		it('rounds to 2 decimal defined values', () => {
			expect(StateComponent.displayValue(12.345)).toBe('12.35');
		});
	});

	describe('static state', () => {
		beforeEach(() => {
			fixture = TestBed.createComponent(StateComponent);
			component = fixture.componentInstance;
			component.current = new State('day',60.0,25,100,90,20);
			component.currentTemperature = 24.3567;
			component.currentHumidity = 64.3456;
			fixture.detectChanges();
		});

		it('should create', () => {
			expect(component).toBeTruthy();
		});
		describe ('no current display',() => {
			it('should display target state',() => {
				const compiled = fixture.debugElement.nativeElement;
				let selection = compiled.querySelectorAll('table tbody td');
				expect(selection.length).toBe(5);
				expect(selection[0].textContent).toContain('25','temperature');
				expect(selection[1].textContent).toContain('60','humidity');
				expect(selection[2].textContent).toContain('100','wind');
				expect(selection[3].textContent).toContain('90','visibleLight');
				expect(selection[4].textContent).toContain('20','uvLight');
			});
		});

		describe ('with current display',() => {
			it('should display target state',() => {
				component.displayCurrent = true;
				fixture.detectChanges();
				const compiled = fixture.debugElement.nativeElement;
				let selection = compiled.querySelectorAll('table tbody td');
				expect(selection.length).toBe(10);
				expect(selection[0].textContent).toContain('25','temperature');
				expect(selection[1].textContent).toContain('60','humidity');
				expect(selection[2].textContent).toContain('100','wind');
				expect(selection[3].textContent).toContain('90','visibleLight');
				expect(selection[4].textContent).toContain('20','uvLight');
				expect(selection[5].textContent).toContain('24.36','current temperature');
				expect(selection[6].textContent).toContain('64.35','current humidity');
				expect(selection[7].textContent).toContain('N.A.','current wind');
				expect(selection[8].textContent).toContain('N.A.','current visibleLight');
				expect(selection[9].textContent).toContain('N.A.','current uvLight');

			});
		});
	});

	describe('transition state', () => {
		beforeEach(() => {
			fixture = TestBed.createComponent(StateComponent);
			component = fixture.componentInstance;
			component.current = new State('day',60.0,24.3567,100,67.34566,18.3456);
			component.end = new State('night',60.0,21,100,0,0);
			component.currentTemperature = 24.3567;
			component.currentHumidity = 64.3456;
			fixture.detectChanges();
		});

		it('should create', () => {
			expect(component).toBeTruthy();
		});
		describe ('no current display',() => {
			it('should display target state',() => {
				const compiled = fixture.debugElement.nativeElement;
				let selection = compiled.querySelectorAll('table tbody td');
				expect(selection.length).toBe(5);
				expect(selection[0].textContent).toContain('24.36 ↦ 21','temperature');
				expect(selection[1].textContent).toContain('60 ↦ 60','humidity');
				expect(selection[2].textContent).toContain('100 ↦ 100','wind');
				expect(selection[3].textContent).toContain('67.35 ↦ 0','visibleLight');
				expect(selection[4].textContent).toContain('18.35 ↦ 0','uvLight');
			});
		});

		describe ('with current display',() => {
			it('should display target state',() => {
				component.displayCurrent = true;
				fixture.detectChanges();
				const compiled = fixture.debugElement.nativeElement;
				let selection = compiled.querySelectorAll('table tbody td');
				expect(selection.length).toBe(10);
				expect(selection[0].textContent).toContain('24.36 ↦ 21','temperature');
				expect(selection[1].textContent).toContain('60 ↦ 60','humidity');
				expect(selection[2].textContent).toContain('100 ↦ 100','wind');
				expect(selection[3].textContent).toContain('67.35 ↦ 0','visibleLight');
				expect(selection[4].textContent).toContain('18.35 ↦ 0','uvLight');
				expect(selection[5].textContent).toContain('24.36','current temperature');
				expect(selection[6].textContent).toContain('64.35','current humidity');
				expect(selection[7].textContent).toContain('N.A.','current wind');
				expect(selection[8].textContent).toContain('N.A.','current visibleLight');
				expect(selection[9].textContent).toContain('N.A.','current uvLight');

			});
		});
	});



});
