import { ComponentFixture, TestBed } from '@angular/core/testing';

import { StatusToggle } from './status-toggle';

describe('StatusToggle', () => {
  let component: StatusToggle;
  let fixture: ComponentFixture<StatusToggle>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [StatusToggle]
    })
    .compileComponents();

    fixture = TestBed.createComponent(StatusToggle);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
