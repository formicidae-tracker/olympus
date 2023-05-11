import { ComponentFixture, TestBed } from '@angular/core/testing';

import { NodeIndexComponent } from './node-index.component';
import { NodeCardComponent } from './node-card/node-card.component';
import { MatCardModule } from '@angular/material/card';



describe('NodeIndexComponent', () => {
  let component: NodeIndexComponent;
  let fixture: ComponentFixture<NodeIndexComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
		declarations: [NodeIndexComponent, NodeCardComponent],
		imports: [ MatCardModule ],
    });
    fixture = TestBed.createComponent(NodeIndexComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
