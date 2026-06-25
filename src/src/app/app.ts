import {Component, OnDestroy, OnInit, signal} from "@angular/core";
import {UpdateService} from "./services/update.service";
import {gitCount} from "./build";

@Component({
    selector: "app-root",
    imports: [],
    templateUrl: "./app.html",
    styleUrl: "./app.scss",
})
export class App implements OnInit, OnDestroy {
    protected readonly title = signal("snooker");

    private readonly _visibilityListener: () => void;

    constructor(private readonly updateSvc: UpdateService) {
        this._visibilityListener = this.visibilityListener.bind(this);
    }

    private visibilityListener(): void {
        this.updateSvc.checkForUpdate().catch(console.error);
    }

    ngOnInit() {
        document.addEventListener("visibilitychange", this._visibilityListener);
        window.addEventListener("focus", this._visibilityListener);
    }

    ngOnDestroy(): void {
        document.removeEventListener(
            "visibilitychange",
            this._visibilityListener,
        );
        window.removeEventListener("focus", this._visibilityListener);
    }

    protected gitCount = gitCount;
}
