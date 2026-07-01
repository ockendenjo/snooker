import {Routes} from "@angular/router";
import {HomePage} from "./home-page/home-page";
import {InfoPage} from "./info-page/info-page";
import {AbvPage} from "./abv-page/abv-page";
import {EdiMapPage} from "./edi-map-page/edi-map-page";
import {CallbackComponent} from "./callback/callback.component";
import {ProfilePage} from "./profile-page/profile-page";
import {AuthService} from "./services/auth.service";
import {requireSignedIn} from "./auth.guards";

export const routes: Routes = [
    {
        path: "",
        component: HomePage,
    },
    {
        path: "callback",
        component: CallbackComponent,
    },
    {
        path: "abvs",
        component: AbvPage,
    },
    {
        path: "info",
        component: InfoPage,
    },
    {
        path: "map",
        component: EdiMapPage,
    },
    {
        path: "profile",
        component: ProfilePage,
        canActivate: [requireSignedIn],
    },
];
