import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ProfitHistory } from './profit-history';

describe('ProfitHistory', () => {
  let component: ProfitHistory;
  let fixture: ComponentFixture<ProfitHistory>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ProfitHistory]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ProfitHistory);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
