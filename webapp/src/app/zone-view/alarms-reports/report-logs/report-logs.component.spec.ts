import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ReportLogsComponent } from './report-logs.component';
import { CoreModule } from 'src/app/core/core.module';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { DatePipe } from '@angular/common';

describe('ReportLogsComponent', () => {
  let component: ReportLogsComponent;
  let fixture: ComponentFixture<ReportLogsComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ReportLogsComponent],
      imports: [CoreModule, BrowserAnimationsModule],
      providers: [DatePipe],
    });
    fixture = TestBed.createComponent(ReportLogsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
