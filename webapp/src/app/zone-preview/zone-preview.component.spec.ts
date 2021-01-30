import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ZonePreviewComponent } from './zone-preview.component';

describe('ZonePreviewComponent', () => {
  let component: ZonePreviewComponent;
  let fixture: ComponentFixture<ZonePreviewComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ ZonePreviewComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ZonePreviewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
