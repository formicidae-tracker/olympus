import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ZoneIndexComponent } from './zone-index.component';
import { ZoneCardComponent } from './zone-card/zone-card.component';
import { MatCardModule } from '@angular/material/card';
import { HttpClientModule } from '@angular/common/http';

describe('ZoneIndexComponent', () => {
  let component: ZoneIndexComponent;
  let fixture: ComponentFixture<ZoneIndexComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ZoneIndexComponent, ZoneCardComponent],
      imports: [MatCardModule, HttpClientModule],
    });
    fixture = TestBed.createComponent(ZoneIndexComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
