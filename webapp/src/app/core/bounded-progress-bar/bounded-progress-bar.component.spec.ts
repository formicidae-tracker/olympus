import { ComponentFixture, TestBed } from '@angular/core/testing';
import { BoundedProgressBarComponent } from './bounded-progress-bar.component';
import { MatProgressBarModule } from '@angular/material/progress-bar';

import { cases } from 'jasmine-parameterized';

describe('BoundedProgressBarComponent', () => {
  let component: BoundedProgressBarComponent;
  let fixture: ComponentFixture<BoundedProgressBarComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [BoundedProgressBarComponent],
      imports: [MatProgressBarModule],
    });
    fixture = TestBed.createComponent(BoundedProgressBarComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should be in query mode by default', () => {
    expect(component.progressMode()).toEqual('query');
  });

  it('should be deterministic if a value is provided', () => {
    component.value = 20.0;
    expect(component.progressMode()).toEqual('determinate');
  });

  cases([
    // min max defaults to [0,100]
    [-1.0, undefined, undefined, 0.0],
    [24.0, undefined, undefined, 24.0],
    [105.0, undefined, undefined, 100.0],
    [-1.0, undefined, 1.0, 0.0],
    [24.0, undefined, 1.0, 100.0],
    [0.24, undefined, 1.0, 24.0],
    [8.0, 18.0, 20.0, 0.0],
    [19.0, 18.0, 20.0, 50.0],
    [21.0, 19.0, 20.0, 100.0],
  ]).it(
    'should compute the right progress from provided bounds',
    ([value, min, max, expected]) => {
      component.minimum = min;
      component.maximum = max;
      component.value = value;
      expect(component.progressValue()).toBeCloseTo(expected, 7);
    }
  );

  cases([
    // min max defaults to [0,100], and on;y primary color
    [undefined, undefined, undefined, 'primary'],
    [-1.0, undefined, undefined, 'primary'],
    [50.0, undefined, undefined, 'primary'],
    [101.0, undefined, undefined, 'primary'],
    // if min is defined, colorize accordingly
    [undefined, 0.0, undefined, 'primary'],
    [-1.0, 0.0, undefined, 'warn'],
    [4.0, 0.0, undefined, 'accent'],
    [50.0, 0.0, undefined, 'primary'],
    [101.0, 0.0, undefined, 'primary'],
    // if max is defined, colorize accordingly
    [undefined, undefined, 100.0, 'primary'],
    [-1.0, undefined, 100.0, 'primary'],
    [40.0, undefined, 100.0, 'primary'],
    [96.0, undefined, 100.0, 'accent'],
    [101.0, undefined, 100.0, 'warn'],
  ]).it(
    'should choose appropriate color from boundaries',
    ([value, min, max, expected]) => {
      component.minimum = min;
      component.maximum = max;
      component.value = value;
      expect(component.progressColor()).toEqual(expected);
    }
  );
});
