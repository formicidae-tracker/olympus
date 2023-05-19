import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ZoneCardComponent } from './zone-card.component';
import { RouterModule } from '@angular/router';
import { CoreModule } from 'src/app/core/core.module';

describe('ZoneCardComponent', () => {
  let component: ZoneCardComponent;
  let fixture: ComponentFixture<ZoneCardComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ZoneCardComponent],
      imports: [CoreModule, RouterModule.forRoot([])],
    });
    fixture = TestBed.createComponent(ZoneCardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
