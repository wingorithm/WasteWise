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
    idx = 0;
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
        this.idx = (this.idx + 1) % 2;
    }
}
