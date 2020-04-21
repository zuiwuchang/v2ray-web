import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { IptablesSaveComponent } from './iptables-save.component';

describe('IptablesSaveComponent', () => {
  let component: IptablesSaveComponent;
  let fixture: ComponentFixture<IptablesSaveComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ IptablesSaveComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(IptablesSaveComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
