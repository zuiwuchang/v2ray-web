import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { V2raySubscriptionComponent } from './v2ray-subscription.component';

describe('V2raySubscriptionComponent', () => {
  let component: V2raySubscriptionComponent;
  let fixture: ComponentFixture<V2raySubscriptionComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ V2raySubscriptionComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(V2raySubscriptionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
