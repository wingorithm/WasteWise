import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DetectPageComponent } from './detect-page.component';

describe('DetectPageComponent', () => {
  let component: DetectPageComponent;
  let fixture: ComponentFixture<DetectPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DetectPageComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(DetectPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
