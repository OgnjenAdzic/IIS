import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ManageRestaurant } from './manage-restaurant';

describe('ManageRestaurant', () => {
  let component: ManageRestaurant;
  let fixture: ComponentFixture<ManageRestaurant>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ManageRestaurant]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ManageRestaurant);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
