import { Component, ViewChild, ViewEncapsulation } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterOutlet } from '@angular/router';
import { IdleScreenComponent } from './idle-screen/idle-screen.component';
import { DetectPageComponent } from './detect-page/detect-page.component';
import { IntroComponent } from './intro/intro.component';
import { ClassifyComponent } from './classify/classify.component';
import { RewardingComponent } from './rewarding/rewarding.component';

@Component({
    selector: 'app-root',
    standalone: true,
    imports: [CommonModule, RouterOutlet, IdleScreenComponent, DetectPageComponent, IntroComponent, ClassifyComponent, RewardingComponent],
    templateUrl: './app.component.html',
    styleUrl: './app.component.css',
    encapsulation: ViewEncapsulation.None
})

export class AppComponent {
    state:string = "reward";

    ngInit() : void {
        this.nextComponent(this.state);
    }

    nextComponent(state : string) : void {
        this.state = state;
    }
}


