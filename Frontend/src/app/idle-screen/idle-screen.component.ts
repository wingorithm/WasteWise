import { HttpHandler } from '@angular/common/http';
import { Component, EventEmitter, Output, ViewEncapsulation } from '@angular/core';

@Component({
  selector: 'app-idle-screen',
  standalone: true,
  imports: [],
  templateUrl: './idle-screen.component.html',
  styleUrl: './idle-screen.component.css',
  encapsulation: ViewEncapsulation.None
})

export class IdleScreenComponent {
    images: string[] = ['Idle - 1.jpg', 'Idle - 2.jpg'];
    idx = 0;
    currentImage!: string;
    handle = setInterval.prototype!;
    
    ngOnInit(): void {
        this.idx = 0
        this.nextImage()
        this.handle = setInterval(() => {this.nextImage()}, 5000)
    }

    ngOnDestroy(): void {
        clearInterval(this.handle);
    }
    
    nextImage() : void {
        this.currentImage = '/assets/page/' + this.images[this.idx];
        console.log(this.currentImage)
        this.idx = (this.idx + 1) % this.images.length;
    }
}
