import { ComponentFixture, TestBed } from '@angular/core/testing';

import { IdleScreenComponent } from './idle-screen.component';

describe('IdleScreenComponent', () => {
  let component: IdleScreenComponent;
  let fixture: ComponentFixture<IdleScreenComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [IdleScreenComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(IdleScreenComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
