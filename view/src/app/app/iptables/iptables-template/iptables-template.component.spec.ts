import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { IptablesTemplateComponent } from './iptables-template.component';

describe('IptablesTemplateComponent', () => {
  let component: IptablesTemplateComponent;
  let fixture: ComponentFixture<IptablesTemplateComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ IptablesTemplateComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(IptablesTemplateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
