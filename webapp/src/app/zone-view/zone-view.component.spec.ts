import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ZoneViewComponent } from './zone-view.component';
import { RouterModule } from '@angular/router';
import { HttpClientModule } from '@angular/common/http';
import { CoreModule } from '../core/core.module';

describe('ZoneViewComponent', () => {
  let component: ZoneViewComponent;
  let fixture: ComponentFixture<ZoneViewComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [RouterModule.forRoot([]), HttpClientModule, CoreModule],
      declarations: [ZoneViewComponent],
    });
    fixture = TestBed.createComponent(ZoneViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
