import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TrackingStatusComponent } from './tracking-status.component';
import { CoreModule } from 'src/app/core/core.module';

describe('TrackingStatusComponent', () => {
  let component: TrackingStatusComponent;
  let fixture: ComponentFixture<TrackingStatusComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [TrackingStatusComponent],
      imports: [CoreModule],
    });
    fixture = TestBed.createComponent(TrackingStatusComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
