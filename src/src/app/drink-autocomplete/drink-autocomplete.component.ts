import {Component, EventEmitter, Input, Output} from "@angular/core";
import {FormControl, ReactiveFormsModule} from "@angular/forms";
import {MatFormFieldModule} from "@angular/material/form-field";
import {MatInputModule} from "@angular/material/input";
import {MatAutocompleteModule} from "@angular/material/autocomplete";
import {Beer} from "../services/beers.service";
import {BeerImageComponent} from "../beer-image/beer-image.component";
import {searchBeers} from "./search";
import {NgClass} from "@angular/common";

@Component({
    selector: "app-drink-autocomplete",
    imports: [
        ReactiveFormsModule,
        MatFormFieldModule,
        MatInputModule,
        MatAutocompleteModule,
        BeerImageComponent,
        NgClass,
    ],
    templateUrl: "./drink-autocomplete.component.html",
    styleUrl: "./drink-autocomplete.component.scss",
})
export class DrinkAutocompleteComponent {
    @Input({required: true}) allBeers: Beer[] = [];
    @Input({required: true}) control!: FormControl;
    @Output() beerSelected = new EventEmitter<Beer>();
    @Output() inputChanged = new EventEmitter<void>();

    public get filteredBeers(): Beer[] {
        const value = this.control.value;
        const search = typeof value === "string" ? value : (value?.name ?? "");
        return searchBeers(this.allBeers, search);
    }

    public getBeerText(v: Beer | string | null): string {
        if (v === null) {
            return "";
        }
        if (typeof v === "object") {
            return v.name;
        }
        return v;
    }

    public getColour(b: Beer): string {
        if (b.abv >= 5.4) {
            return "black";
        }
        if (b.abv >= 5.1) {
            return "pink";
        }
        if (b.abv >= 4.8) {
            return "blue";
        }
        if (b.abv >= 4.5) {
            return "brown";
        }
        if (b.abv >= 4.2) {
            return "green";
        }
        if (b.abv >= 3.9) {
            return "yellow";
        }
        return "red";
    }
}
