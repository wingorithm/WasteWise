import { Component, ViewChild, ViewEncapsulation } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterOutlet } from '@angular/router';
import { IdleScreenComponent } from './idle-screen/idle-screen.component';
import { DetectPageComponent } from './detect-page/detect-page.component';

@Component({
    selector: 'app-root',
    standalone: true,
    imports: [CommonModule, RouterOutlet, IdleScreenComponent, DetectPageComponent],
    templateUrl: './app.component.html',
    styleUrl: './app.component.css',
    encapsulation: ViewEncapsulation.None
})

export class AppComponent {
    @ViewChild('idle', {static:false}) idleScreen!: IdleScreenComponent
    @ViewChild('main', {static:false}) detectPage!: DetectPageComponent


    state:string = "idle";
    mapping: mapperDict = {};

    ngAfterViewInit() : void {
        this.mapping = {
            "idle": this.idleScreen,
            "main": this.detectPage
        }
        this.nextComponent(this.state);
    }

    nextComponent(state : string) : void {
        this.state = state;
        this.mapping[state].init();
    }
}

interface mapperDict {
    [index:string]: BaseComponent;
}

