import { ComponentFixture, TestBed } from '@angular/core/testing';

import { LogIndexComponent } from './log-index.component';
import { HttpClientModule } from '@angular/common/http';
import { MatCardModule } from '@angular/material/card';
import { MatExpansionModule } from '@angular/material/expansion';

describe('LogIndexComponent', () => {
  let component: LogIndexComponent;
  let fixture: ComponentFixture<LogIndexComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [LogIndexComponent],
      imports: [HttpClientModule, MatCardModule, MatExpansionModule],
    });
    fixture = TestBed.createComponent(LogIndexComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
