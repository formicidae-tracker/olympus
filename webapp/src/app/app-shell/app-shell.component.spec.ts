import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AppShellComponent } from './app-shell.component';

import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

describe('AppShellComponent', () => {
  let component: AppShellComponent;
  let fixture: ComponentFixture<AppShellComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [AppShellComponent],
      imports: [MatProgressSpinnerModule],
    });
    fixture = TestBed.createComponent(AppShellComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
