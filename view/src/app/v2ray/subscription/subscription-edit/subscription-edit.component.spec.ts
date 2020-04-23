import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SubscriptionEditComponent } from './subscription-edit.component';

describe('SubscriptionEditComponent', () => {
  let component: SubscriptionEditComponent;
  let fixture: ComponentFixture<SubscriptionEditComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SubscriptionEditComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SubscriptionEditComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
