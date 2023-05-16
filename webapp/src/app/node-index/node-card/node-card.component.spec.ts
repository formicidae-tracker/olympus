import { ComponentFixture, TestBed } from '@angular/core/testing';

import { NodeCardComponent } from './node-card.component';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';

describe('NodeCardComponent', () => {
  let component: NodeCardComponent;
  let fixture: ComponentFixture<NodeCardComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [NodeCardComponent],
      imports: [MatCardModule, MatIconModule],
    });
    fixture = TestBed.createComponent(NodeCardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
