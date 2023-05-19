import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AlarmsReportsComponent } from './alarms-reports.component';
import { CoreModule } from 'src/app/core/core.module';

describe('AlarmsReportsComponent', () => {
  let component: AlarmsReportsComponent;
  let fixture: ComponentFixture<AlarmsReportsComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [AlarmsReportsComponent],
      imports: [CoreModule],
    });
    fixture = TestBed.createComponent(AlarmsReportsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
