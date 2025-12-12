import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AddItemForm } from './add-item-form';

describe('AddItemForm', () => {
  let component: AddItemForm;
  let fixture: ComponentFixture<AddItemForm>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AddItemForm]
    })
    .compileComponents();

    fixture = TestBed.createComponent(AddItemForm);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
