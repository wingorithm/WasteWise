import { Component, EventEmitter, Output, ViewEncapsulation } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-intro',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './intro.component.html',
  styleUrl: './intro.component.css'
})

export class IntroComponent {
    @Output() switch = new EventEmitter<void>();
    idx = 1
    isTransition = true
    text = "Hallo teman - teman, perkenalkan nama aku Mila."

    ngOnInit(): void {
        this.idx = 1
        this.isTransition = true

        setTimeout(() => {
            this.idx = 2;
            setTimeout(() => {
                this.isTransition = false;
                this.enterIntro();
            }, 1500)
        }, 1500)
        // setTimeout(() => {
        //     this.switch.emit();
        // }, 3000);
    }

    enterIntro(): void {
        setTimeout(() => {
            this.text = "Selamat datang di WasteWise";
            setTimeout(() => {
                this.switch.emit();
            }, 2500)
        }, 3000)
    }
}
