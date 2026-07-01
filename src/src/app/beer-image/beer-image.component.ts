import {Component, Input} from "@angular/core";
import {Beer} from "../services/beers.service";
import {NgOptimizedImage} from "@angular/common";
import {ApiService} from "../services/api.service";

@Component({
    selector: "app-beer-image",
    standalone: true,
    templateUrl: "./beer-image.component.html",
    imports: [NgOptimizedImage],
})
export class BeerImageComponent {
    @Input({required: true}) beer!: Beer;
    @Input() height!: number;

    constructor(private readonly apiSvc: ApiService) {}

    public handleError(evt: ErrorEvent) {
        if (evt.target) {
            const target = evt.target as HTMLImageElement;
            const src = target.src;
            if (src) {
                const body = {
                    src,
                    name: this.beer.name,
                    brewery: this.beer.brewery,
                };
                this.apiSvc.send("reportBrokenImg", body).catch(undefined);
            }
        }
    }
}
