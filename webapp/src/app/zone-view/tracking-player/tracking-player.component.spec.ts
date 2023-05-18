import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TrackingPlayerComponent } from './tracking-player.component';

describe('TrackingPlayerComponent', () => {
  let component: TrackingPlayerComponent;
  let fixture: ComponentFixture<TrackingPlayerComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [TrackingPlayerComponent]
    });
    fixture = TestBed.createComponent(TrackingPlayerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
