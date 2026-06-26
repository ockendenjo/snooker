import {Component, OnInit} from "@angular/core";
import {Pub, PubsService} from "../services/pubs.service";
import {renderMap} from "./map";

@Component({
    selector: "app-edi-map-page",
    imports: [],
    templateUrl: "./edi-map-page.html",
    styleUrl: "./edi-map-page.scss",
})
export class EdiMapPage implements OnInit {
    constructor(private readonly pubsSvc: PubsService) {}

    ngOnInit() {
        this.pubsSvc.getPubs().then((pubs: Pub[]) => {
            renderMap(pubs, () => {});
        });
    }
}
