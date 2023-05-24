import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ClimateStatusComponent } from './climate-status.component';
import { CoreModule } from 'src/app/core/core.module';

describe('ClimateStatusComponent', () => {
  let component: ClimateStatusComponent;
  let fixture: ComponentFixture<ClimateStatusComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ClimateStatusComponent],
      imports: [CoreModule],
    });
    fixture = TestBed.createComponent(ClimateStatusComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
