import {Injectable} from "@angular/core";
import {SwUpdate, VersionReadyEvent} from "@angular/service-worker";
import {filter} from "rxjs";

@Injectable({providedIn: "root"})
export class UpdateService {
    constructor(private readonly swUpdate: SwUpdate) {
        if (!swUpdate.isEnabled) return;

        swUpdate.versionUpdates
            .pipe(
                filter(
                    (e): e is VersionReadyEvent => e.type === "VERSION_READY",
                ),
            )
            .subscribe(() => {
                swUpdate
                    .activateUpdate()
                    .then(() => document.location.reload());
            });

        swUpdate.unrecoverable.subscribe(() => {
            document.location.reload();
        });
    }

    checkForUpdate(): Promise<boolean> {
        if (!this.swUpdate.isEnabled) return Promise.resolve(false);
        return this.swUpdate.checkForUpdate();
    }
}
