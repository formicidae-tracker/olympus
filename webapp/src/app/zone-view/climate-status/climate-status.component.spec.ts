import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ClimateStatusComponent } from './climate-status.component';

describe('ClimateStatusComponent', () => {
  let component: ClimateStatusComponent;
  let fixture: ComponentFixture<ClimateStatusComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ClimateStatusComponent],
    });
    fixture = TestBed.createComponent(ClimateStatusComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
