import {ChangeDetectorRef, Component, OnInit, signal} from "@angular/core";
import {ActivatedRoute} from "@angular/router";
import {
    FormBuilder,
    FormGroup,
    ReactiveFormsModule,
    Validators,
} from "@angular/forms";
import {MatButtonModule} from "@angular/material/button";
import {MatFormFieldModule, MatLabel} from "@angular/material/form-field";
import {MatInputModule} from "@angular/material/input";
import {DatePipe} from "@angular/common";
import {MatOption} from "@angular/material/core";
import {MatSelectModule} from "@angular/material/select";
import {MatRadioModule} from "@angular/material/radio";
import {MatAutocompleteModule} from "@angular/material/autocomplete";
import {MatCheckboxModule} from "@angular/material/checkbox";
import {DrinkService, NewDrink} from "../services/drink.service";
import {Pub, PubsService} from "../services/pubs.service";
import {Beer, BeersService} from "../services/beers.service";
import {RangeService} from "../services/range.service";
import {environment} from "../../environments/environment.dev";
import {DrinkAutocompleteComponent} from "../drink-autocomplete/drink-autocomplete.component";
import {RatingControl} from "../rating-control/rating-control";

@Component({
    selector: "app-log-drink-page",
    imports: [
        ReactiveFormsModule,
        MatLabel,
        MatOption,
        MatSelectModule,
        MatFormFieldModule,
        MatInputModule,
        MatButtonModule,
        MatRadioModule,
        DatePipe,
        MatAutocompleteModule,
        MatCheckboxModule,
        DrinkAutocompleteComponent,
        RatingControl,
    ],
    templateUrl: "./log-drink-page.html",
    styleUrl: "./log-drink-page.scss",
})
export class LogDrinkPage implements OnInit {
    public tooEarly = false;
    public tooLate = false;
    public debugUI = false;
    public pageState = signal(PageState.When);
    public errorStr = this.debugUI ? "Something went wrong" : "";

    public day0 = 0;
    public day1 = 0;

    protected pubIDFromRoute: number | null = null;

    public whenForm: FormGroup;
    public pubForm: FormGroup;
    public drinkForm: FormGroup;
    public extraForm: FormGroup;

    private selectedPubID: number | null = null;
    private selectedPubs: Pub[] = [];
    public filteredPubs: Pub[] = [];
    public allBeers: Beer[] = [];

    constructor(
        private readonly drinkSvc: DrinkService,
        private readonly pubsSvc: PubsService,
        private readonly beersSvc: BeersService,
        private readonly route: ActivatedRoute,
        private readonly fb: FormBuilder,
        private readonly cdr: ChangeDetectorRef,
        rangeSvc: RangeService,
    ) {
        this.tooEarly = rangeSvc.isTooEarly();
        this.tooLate = rangeSvc.isTooLate();

        const initialDaySelect = this.setupDateSelect();

        this.whenForm = this.fb.group({
            selectWhen: ["now"],
            day_select: [
                {value: initialDaySelect, disabled: true},
                Validators.required,
            ],
            time: [{value: "", disabled: true}, Validators.required],
        });

        this.drinkForm = this.fb.group({
            name: ["", Validators.required],
            brewery: ["", Validators.required],
            untappdID: [null as number | null],
            abv: [null as number | null, Validators.required],
        });

        this.pubForm = this.fb.group({
            venue: ["", Validators.required],
            drunkWith: ["", [Validators.required, Validators.maxLength(100)]],
        });

        this.pubForm.get("venue")!.valueChanges.subscribe(() => {
            this.updateFilteredPubs();
        });

        this.extraForm = this.fb.group({
            notes: ["", Validators.maxLength(200)],
            pubRating: [null as number | null],
            price: [null as number | null, Validators.min(0)],
        });

        this.whenForm.get("selectWhen")!.valueChanges.subscribe((val) => {
            const dayCtrl = this.whenForm.get("day_select")!;
            const timeCtrl = this.whenForm.get("time")!;
            if (val === "past") {
                dayCtrl.enable();
                timeCtrl.enable();
            } else {
                dayCtrl.disable();
                timeCtrl.disable();
            }
            this.cdr.markForCheck();
        });
    }

    private setupDateSelect(): string {
        const dayInMillis = 24 * 60 * 60 * 1_000;
        const now = Date.now();
        this.day0 = now;
        let initialDaySelect = "day0";
        if (environment.endDate && this.day0 > environment.endDate.getTime()) {
            this.day0 = 0;
            initialDaySelect = "day1";
        }
        this.day1 = now - dayInMillis;
        if (
            environment.startDate &&
            this.day1 < environment.startDate.getTime()
        ) {
            this.day1 = 0;
        }
        return initialDaySelect;
    }

    ngOnInit(): void {
        const pubIDParam = this.route.snapshot.queryParamMap.get("pubID");
        if (pubIDParam) {
            this.pubIDFromRoute = parseInt(pubIDParam, 10);
            this.transitionToPubDetails();
        }
    }

    private updateFilteredPubs(): void {
        const venue = this.pubForm.get("venue")!.value ?? "";
        if (!venue) {
            this.filteredPubs = this.selectedPubs;
            return;
        }
        const lower = venue.toLowerCase();
        this.filteredPubs = this.selectedPubs.filter((s) =>
            s.name.toLowerCase().includes(lower),
        );
    }

    public onPubSelected(pub: Pub | string): void {
        if (typeof pub === "object") {
            this.selectedPubID = pub.camraID;
            this.pubForm.patchValue({venue: pub.name});
        } else {
            this.selectedPubID = null;
        }
    }

    public onBeerSelected(beer: Beer): void {
        this.drinkForm.patchValue({
            name: beer.name,
            brewery: beer.brewery,
            untappdID: beer.untappd,
            abv: beer.abv,
        });
    }

    public onNameChange(): void {
        this.drinkForm.patchValue({brewery: "", untappdID: null, abv: null});
    }

    public logDrink(): void {
        const p = this.pubForm.value;
        const d = this.drinkForm.value;
        const e = this.extraForm.value;

        const drink: NewDrink = {
            timestamp: this.getDrinkTimestamp(),
            pubName: p.venue,
            camraID: this.selectedPubID || undefined,
            drinkName: d.name,
            brewery: d.brewery,
            abv: d.abv,
            untappdID: d.untappdID || undefined,
            with: p.drunkWith,
            notes: e.notes || undefined,
        };

        this.pageState.set(PageState.Saving);
        this.drinkSvc
            .logDrink(drink)
            .then(() => {
                this.pageState.set(PageState.Saved);
            })
            .catch((e) => {
                this.errorStr = e;
                this.pageState.set(PageState.Error);
            });
    }

    public transitionToDrinkDetails(): void {
        if (this.allBeers.length) {
            this.pageState.set(PageState.DrinkDetails);
            return;
        }

        this.pageState.set(PageState.LoadingDrinks);

        this.beersSvc
            .loadAll()
            .then((beers) => {
                this.allBeers = beers;
                this.pageState.set(PageState.DrinkDetails);
            })
            .catch(() => {
                this.errorStr = "Failed to load drink details";
                this.pageState.set(PageState.Error);
            });
    }

    public transitionToPubDetails(): void {
        this.pageState.set(PageState.LoadingPubs);

        this.pubsSvc
            .getPubs()
            .then((pubs) => {
                this.selectedPubs = pubs;
                this.updateFilteredPubs();

                if (this.pubIDFromRoute !== null) {
                    const match = this.selectedPubs.find(
                        (p) => p.camraID === this.pubIDFromRoute,
                    );
                    if (match) {
                        this.onPubSelected(match);
                    }
                    this.pubIDFromRoute = null;
                }

                this.pageState.set(PageState.PubDetails);
            })
            .catch(() => {
                this.errorStr = "Failed to load pub selections";
                this.pageState.set(PageState.Error);
            });
    }

    public transitionToExtraDetails(): void {
        this.pageState.set(PageState.ExtraDetails);
    }

    public goBackWhen(): void {
        this.pageState.set(PageState.When);
    }

    public goBackPub(): void {
        this.pageState.set(PageState.PubDetails);
    }

    public goBackDrink(): void {
        this.pageState.set(PageState.DrinkDetails);
    }

    public goAgain(): void {
        const initialDaySelect = this.setupDateSelect();
        this.whenForm.reset({
            selectWhen: "now",
            day_select: initialDaySelect,
            time: "",
        });
        this.pubForm.reset({
            venue: "",
            drunkWith: "",
        });
        this.drinkForm.reset({
            name: "",
            brewery: "",
            untappdID: null,
            abv: null,
        });
        this.extraForm.reset({
            notes: "",
            pubRating: null,
            price: null,
        });
        this.selectedPubID = null;
        this.selectedPubs = [];
        this.filteredPubs = [];
        this.errorStr = "";
        this.pageState.set(PageState.When);
    }

    private getDrinkTimestamp(): string | undefined {
        const {selectWhen, day_select, time} = this.whenForm.getRawValue();
        if (selectWhen === "now") {
            return undefined;
        }

        let t = Date.now();
        if (day_select === "day1") {
            t -= 24 * 60 * 60 * 1_000;
        }
        if (day_select === "day2") {
            t -= 2 * 24 * 60 * 60 * 1_000;
        }

        const d = new Date(t);
        const parts = (time as string).split(":");
        if (parts.length === 2) {
            d.setHours(parseInt(parts[0]), parseInt(parts[1]));
        }

        return d.toISOString();
    }
}

enum PageState {
    When = "When",
    LoadingPubs = "LoadingPubs",
    PubDetails = "PubDetails",
    LoadingDrinks = "LoadingDrinks",
    DrinkDetails = "DrinkDetails",
    ExtraDetails = "ExtraDetails",
    Saving = "Saving",
    Error = "Error",
    Saved = "Saved",
}
