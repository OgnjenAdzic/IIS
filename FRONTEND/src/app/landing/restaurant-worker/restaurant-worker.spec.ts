import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RestaurantWorker } from './restaurant-worker';

describe('RestaurantWorker', () => {
  let component: RestaurantWorker;
  let fixture: ComponentFixture<RestaurantWorker>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RestaurantWorker]
    })
    .compileComponents();

    fixture = TestBed.createComponent(RestaurantWorker);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
