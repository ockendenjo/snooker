import {ChangeDetectorRef, Component, OnInit} from "@angular/core";
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
    ],
    templateUrl: "./log-drink-page.html",
    styleUrl: "./log-drink-page.scss",
})
export class LogDrinkPage implements OnInit {
    public tooEarly = false;
    public tooLate = false;
    public debugUI = false;
    public pageState: PageState = PageState.Ready;
    public errorStr = this.debugUI ? "Something went wrong" : "";

    public day0 = 0;
    public day1 = 0;

    protected pubIDFromRoute: number | null = null;

    public whenForm: FormGroup;
    public detailsForm: FormGroup;

    private selectedPubID: number | null = null;
    private selectedPubs: Pub[] = [];
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

        this.detailsForm = this.fb.group({
            venue: ["", Validators.required],
            name: ["", Validators.required],
            brewery: ["", Validators.required],
            untappdID: [null as number | null],
            drunkWith: ["", [Validators.required, Validators.maxLength(100)]],
            endOfWord: [false],
            notInWord: [false],
            notes: ["", Validators.maxLength(200)],
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
            this.goDetails();
        }
    }

    public get filteredPubs(): Pub[] {
        const venue = this.detailsForm.get("venue")!.value ?? "";
        if (!venue) {
            return this.selectedPubs;
        }
        const lower = venue.toLowerCase();
        return this.selectedPubs.filter((s) => {
            return s.name.toLowerCase().includes(lower);
        });
    }

    public onPubSelected(pub: Pub): void {
        this.selectedPubID = pub.camraID;
        this.detailsForm.patchValue({venue: pub.name});
    }

    public onBeerSelected(beer: Beer): void {
        this.detailsForm.patchValue({
            name: beer.name,
            brewery: beer.brewery,
            untappdID: beer.untappd,
        });
    }

    public onNameChange(): void {
        this.detailsForm.patchValue({brewery: "", untappdID: null});
    }

    public cid = "";

    public logDrink(): void {
        if (this.selectedPubID === null) {
            this.errorStr = "Please select a valid pub from the list";
            this.pageState = PageState.Error;
            return;
        }
        const v = this.detailsForm.value;
        const drink: NewDrink = {
            pubID: this.selectedPubID,
            name: v.name,
            brewery: v.brewery,
            untappdID: v.untappdID || undefined,
            with: v.drunkWith,
            timestamp: this.getDrinkTimestamp(),
            endOfWord: v.endOfWord,
            notInWord: v.notInWord,
            notes: v.notes || undefined,
        };
        this.cid = "";
        this.pageState = PageState.Saving;
        this.drinkSvc
            .logDrink(drink)
            .then((ld) => {
                this.pageState = PageState.Saved;
                this.cid = ld.cid;
                this.cdr.markForCheck();
            })
            .catch((e) => {
                this.errorStr = e;
                this.pageState = PageState.Error;
                this.cdr.markForCheck();
            });
    }

    public goDetails(): void {
        this.pageState = PageState.Loading;

        const loadBeers = this.beersSvc.loadAll();
        const loadPubs = this.pubsSvc.getPubs();

        Promise.all([loadBeers, loadPubs])
            .then(([beers, pubs]: [Beer[], Pub[]]) => {
                this.allBeers = beers;
                this.selectedPubs = pubs;

                if (this.pubIDFromRoute !== null) {
                    const match = this.selectedPubs.find(
                        (p) => p.camraID === this.pubIDFromRoute,
                    );
                    if (match) {
                        this.onPubSelected(match);
                    }
                    this.pubIDFromRoute = null;
                }

                this.pageState = PageState.Details;
                this.cdr.markForCheck();
            })
            .catch(() => {
                this.errorStr = "Failed to load pub selections";
                this.pageState = PageState.Error;
                this.cdr.markForCheck();
            });
    }

    private getSelectionTimestamp(): string | undefined {
        if (this.whenForm.get("selectWhen")!.value === "now") {
            return undefined;
        }
        return this.getDrinkTimestamp();
    }

    public goBack(): void {
        this.pageState = PageState.Ready;
    }

    public goAgain(): void {
        const initialDaySelect = this.setupDateSelect();
        this.whenForm.reset({
            selectWhen: "now",
            day_select: initialDaySelect,
            time: "",
        });
        this.detailsForm.reset({
            venue: "",
            name: "",
            brewery: "",
            untappdID: null,
            drunkWith: "",
            endOfWord: false,
            notInWord: false,
            notes: "",
        });
        this.selectedPubID = null;
        this.errorStr = "";
        this.pageState = PageState.Ready;
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
    Ready = "Ready",
    Loading = "Loading",
    Details = "Details",
    Saving = "Saving",
    Error = "Error",
    Saved = "Saved",
}
