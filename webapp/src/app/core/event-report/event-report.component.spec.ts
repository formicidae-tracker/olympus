import { ComponentFixture, TestBed } from '@angular/core/testing';

import { EventReportComponent } from './event-report.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { DatePipe } from '@angular/common';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatTableModule } from '@angular/material/table';

describe('EventReportComponent', () => {
  let component: EventReportComponent;
  let fixture: ComponentFixture<EventReportComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [EventReportComponent],
      imports: [MatTableModule, MatPaginatorModule, BrowserAnimationsModule],
      providers: [DatePipe],
    });
    fixture = TestBed.createComponent(EventReportComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
