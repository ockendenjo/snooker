import {Component, forwardRef} from "@angular/core";
import {ControlValueAccessor, NG_VALUE_ACCESSOR} from "@angular/forms";
import {NgClass} from "@angular/common";

@Component({
    selector: "app-rating-control",
    imports: [NgClass],
    templateUrl: "./rating-control.html",
    styleUrl: "./rating-control.scss",
    providers: [
        {
            provide: NG_VALUE_ACCESSOR,
            useExisting: forwardRef(() => RatingControl),
            multi: true,
        },
    ],
})
export class RatingControl implements ControlValueAccessor {
    public value: number | null = null;
    public isDisabled = false;
    public readonly ratings = [1, 2, 3, 4, 5];

    private onChange: (value: number | null) => void = () => {};
    private onTouched: () => void = () => {};

    public select(rating: number): void {
        this.value = rating;
        this.onChange(rating);
        this.onTouched();
    }

    writeValue(value: number | null): void {
        this.value = value;
    }

    registerOnChange(fn: (value: number | null) => void): void {
        this.onChange = fn;
    }

    registerOnTouched(fn: () => void): void {
        this.onTouched = fn;
    }

    setDisabledState(isDisabled: boolean): void {
        this.isDisabled = isDisabled;
    }
}
