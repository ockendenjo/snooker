import {AfterViewInit, Component, OnInit} from "@angular/core";
import {Pub, PubsService} from "../services/pubs.service";
import {renderMap} from "./map";
import {Router} from "@angular/router";

@Component({
    selector: "app-edi-map-page",
    imports: [],
    templateUrl: "./edi-map-page.html",
    styleUrls: ["./edi-map-page.scss", "../popup.scss"],
})
export class EdiMapPage implements OnInit, AfterViewInit {
    private initFinished = false;
    private viewInitFinished = false;

    private pubs: Pub[] = [];

    constructor(
        private readonly pubsSvc: PubsService,
        private readonly router: Router,
    ) {}

    ngOnInit() {
        this.pubsSvc.getPubs().then((pubs: Pub[]) => {
            this.pubs = pubs;
            this.initFinished = true;
            this.checkFunc();
        });
    }

    ngAfterViewInit() {
        setTimeout(() => {
            this.viewInitFinished = true;
            this.checkFunc();
        }, 0);
    }

    private checkFunc() {
        if (this.initFinished && this.viewInitFinished) {
            renderMap(this.pubs, (pubCamraID) => {
                this.navigateTo(pubCamraID);
            });
        }
    }

    private navigateTo(pubID: number) {
        this.router.navigate(["/log"], {queryParams: {pubID}});
    }
}
