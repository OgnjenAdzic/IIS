import { TestBed } from '@angular/core/testing';

import { Stakeholders } from './stakeholders';

describe('Stakeholders', () => {
  let service: Stakeholders;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(Stakeholders);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
