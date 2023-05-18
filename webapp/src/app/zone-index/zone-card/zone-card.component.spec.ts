import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ZoneCardComponent } from './zone-card.component';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { RouterModule } from '@angular/router';

describe('ZoneCardComponent', () => {
  let component: ZoneCardComponent;
  let fixture: ComponentFixture<ZoneCardComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ZoneCardComponent],
      imports: [MatCardModule, MatIconModule, RouterModule.forRoot([])],
    });
    fixture = TestBed.createComponent(ZoneCardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
