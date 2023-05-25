import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SnackNetworkOfflineComponent } from './snack-network-offline.component';

describe('SnackNetworkOfflineComponent', () => {
  let component: SnackNetworkOfflineComponent;
  let fixture: ComponentFixture<SnackNetworkOfflineComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [SnackNetworkOfflineComponent],
    });
    fixture = TestBed.createComponent(SnackNetworkOfflineComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
