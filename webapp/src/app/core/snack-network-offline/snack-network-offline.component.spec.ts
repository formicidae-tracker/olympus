import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SnackNetworkOfflineComponent } from './snack-network-offline.component';
import { MatSnackBarRef } from '@angular/material/snack-bar';
import { CoreModule } from '../core.module';

describe('SnackNetworkOfflineComponent', () => {
  let component: SnackNetworkOfflineComponent;
  let fixture: ComponentFixture<SnackNetworkOfflineComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [SnackNetworkOfflineComponent],
      imports: [CoreModule],
      providers: [{ provide: MatSnackBarRef, useValue: {} }],
    });
    fixture = TestBed.createComponent(SnackNetworkOfflineComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
