import { ComponentFixture, TestBed } from '@angular/core/testing';

import { VersionComponent } from './version.component';
import { HttpClientModule } from '@angular/common/http';
import { MatListModule } from '@angular/material/list';

describe('VersionComponent', () => {
  let component: VersionComponent;
  let fixture: ComponentFixture<VersionComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [VersionComponent],
      imports: [HttpClientModule, MatListModule],
    });
    fixture = TestBed.createComponent(VersionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
