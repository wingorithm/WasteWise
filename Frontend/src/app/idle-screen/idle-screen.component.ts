import { Component, ViewEncapsulation } from '@angular/core';

@Component({
  selector: 'app-idle-screen',
  standalone: true,
  imports: [],
  templateUrl: './idle-screen.component.html',
  styleUrl: './idle-screen.component.css',
  encapsulation: ViewEncapsulation.None
})

export class IdleScreenComponent {
    images: string[] = ['Idle1.png', 'Idle2.png'];
    idx = 0;
    currentImage: string = '';

    ngOnInit(): void {
        this.nextImage()
        setInterval(() => {this.nextImage()}, 5000)
    }
    
    nextImage() : void {
        this.currentImage = '/assets/' + this.images[this.idx];
        console.log(this.currentImage)
        this.idx = (this.idx + 1) % this.images.length;
    }
}
