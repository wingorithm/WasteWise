import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-classify',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './classify.component.html',
  styleUrl: './classify.component.css'
})

export class ClassifyComponent {
    showAll = true
    @Input() show = 1
    state = 0

    ngOnInit() : void {
        setTimeout(() => {
            this.state = 1;
            setTimeout(() => {
                this.state = 2
                setTimeout(() => {
                    this.showAll = false;
                    setTimeout(() => {
                        this.state = 3;
                    }, 4000)
                }, 2500)
            }, 3000)
        }, 3000)
    }
}
