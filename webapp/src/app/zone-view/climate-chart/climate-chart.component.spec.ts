import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ClimateChartComponent } from './climate-chart.component';

describe('ClimateChartComponent', () => {
  let component: ClimateChartComponent;
  let fixture: ComponentFixture<ClimateChartComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ClimateChartComponent]
    });
    fixture = TestBed.createComponent(ClimateChartComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
