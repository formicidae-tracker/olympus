import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ZoneStatusComponent } from './zone-status.component';

describe('ZoneStatusComponent', () => {
  let component: ZoneStatusComponent;
  let fixture: ComponentFixture<ZoneStatusComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ZoneStatusComponent]
    });
    fixture = TestBed.createComponent(ZoneStatusComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
