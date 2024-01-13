import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';

@Component({
  selector: 'app-rewarding',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './rewarding.component.html',
  styleUrl: './rewarding.component.css'
})
export class RewardingComponent {
    state = 1

    ngOnInit(): void {
        this.state = 1;
        
        setTimeout(() => {
            this.state = 2;
            setTimeout(() => {
                this.state = 3;
            }, 4000)
        }, 4000)        
    }
}
