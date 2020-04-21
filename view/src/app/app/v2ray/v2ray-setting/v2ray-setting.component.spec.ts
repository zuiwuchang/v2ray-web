import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { V2raySettingComponent } from './v2ray-setting.component';

describe('V2raySettingComponent', () => {
  let component: V2raySettingComponent;
  let fixture: ComponentFixture<V2raySettingComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ V2raySettingComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(V2raySettingComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
