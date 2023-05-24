import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ZoneNotificationButtonComponent } from './zone-notification-button.component';
import { CoreModule } from '../core.module';

describe('ZoneNotificationButtonComponent', () => {
  let component: ZoneNotificationButtonComponent;
  let fixture: ComponentFixture<ZoneNotificationButtonComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ZoneNotificationButtonComponent],
      imports: [CoreModule],
    });
    fixture = TestBed.createComponent(ZoneNotificationButtonComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
