import {Injectable} from "@angular/core";
import {environment} from "../../environments/environment.dev";

@Injectable({
    providedIn: "any",
})
export class RangeService {
    isTooEarly(): boolean {
        if (!environment.startDate) {
            return false;
        }
        return new Date() < environment.startDate;
    }

    isTooLate(): boolean {
        if (!environment.endDate) {
            return false;
        }
        const loggingEndDate = new Date(
            environment.endDate.getTime() + 24 * 60 * 60 * 1000,
        );
        return new Date() > loggingEndDate;
    }
}
