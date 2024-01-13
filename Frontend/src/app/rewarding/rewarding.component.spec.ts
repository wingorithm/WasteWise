import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RewardingComponent } from './rewarding.component';

describe('RewardingComponent', () => {
  let component: RewardingComponent;
  let fixture: ComponentFixture<RewardingComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RewardingComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(RewardingComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
