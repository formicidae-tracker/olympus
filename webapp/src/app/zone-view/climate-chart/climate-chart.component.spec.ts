import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ClimateChartComponent } from './climate-chart.component';
import { CoreModule } from 'src/app/core/core.module';
import { NgxEchartsModule } from 'ngx-echarts';
import { HttpClientModule } from '@angular/common/http';

describe('ClimateChartComponent', () => {
  let component: ClimateChartComponent;
  let fixture: ComponentFixture<ClimateChartComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ClimateChartComponent],
      imports: [
        CoreModule,
        NgxEchartsModule.forRoot({ echarts: () => import('echarts') }),
        HttpClientModule,
      ],
    });
    fixture = TestBed.createComponent(ClimateChartComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
