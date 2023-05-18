import { ComponentFixture, TestBed } from '@angular/core/testing';

import { NavToolbarComponent } from './nav-toolbar.component';
import { CoreModule } from '../core/core.module';
import { RouterModule } from '@angular/router';

describe('NavToolbarComponent', () => {
  let component: NavToolbarComponent;
  let fixture: ComponentFixture<NavToolbarComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [NavToolbarComponent],
      imports: [CoreModule, RouterModule.forRoot([])],
    });
    fixture = TestBed.createComponent(NavToolbarComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
