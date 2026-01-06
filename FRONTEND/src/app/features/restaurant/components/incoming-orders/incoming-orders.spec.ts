import { ComponentFixture, TestBed } from '@angular/core/testing';

import { IncomingOrders } from './incoming-orders';

describe('IncomingOrders', () => {
  let component: IncomingOrders;
  let fixture: ComponentFixture<IncomingOrders>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [IncomingOrders]
    })
    .compileComponents();

    fixture = TestBed.createComponent(IncomingOrders);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
