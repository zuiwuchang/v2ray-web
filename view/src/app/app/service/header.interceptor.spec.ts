import { TestBed } from '@angular/core/testing';

import { HeaderInterceptor } from './header.interceptor';

describe('HeaderInterceptor', () => {
  beforeEach(() => TestBed.configureTestingModule({
    providers: [
      HeaderInterceptor
      ]
  }));

  it('should be created', () => {
    const interceptor: HeaderInterceptor = TestBed.inject(HeaderInterceptor);
    expect(interceptor).toBeTruthy();
  });
});
