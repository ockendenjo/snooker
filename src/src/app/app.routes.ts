import {Routes} from "@angular/router";
import {HomePage} from "./home-page/home-page";
import {AbvPage} from "./abv-page/abv-page";

export const routes: Routes = [
    {
        path: "",
        component: HomePage,
    },
    {
        path: "abvs",
        component: AbvPage,
    },
];
