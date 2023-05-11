import { ComponentFixture, TestBed } from '@angular/core/testing';

import { NodeCardComponent } from './node-card.component';

describe('NodeCardComponent', () => {
  let component: NodeCardComponent;
  let fixture: ComponentFixture<NodeCardComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [NodeCardComponent]
    });
    fixture = TestBed.createComponent(NodeCardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
