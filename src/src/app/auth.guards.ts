import {SessionService} from "./services/session.service";
import {AuthService} from "./services/auth.service";
import {CanActivateFn, RedirectCommand, Router} from "@angular/router";
import {inject} from "@angular/core";

export const requireSignedIn: CanActivateFn = () => {
    const authService = inject(AuthService);
    return authService
        .isAuthenticated()
        .then((isAuth) => {
            if (isAuth) {
                return true;
            }
            authService.login();
            return false;
        })
        .catch(() => {
            authService.login();
            return false;
        });
};

export const requireDisplayName: CanActivateFn = () => {
    const router = inject(Router);
    return inject(SessionService)
        .getSessionData()
        .then((sd) => {
            if (sd.displayName.length) {
                return true;
            }
            const rdrTree = router.parseUrl("profile");
            return new RedirectCommand(rdrTree, {});
        })
        .catch(() => {
            const rdrTree = router.parseUrl("profile");
            return new RedirectCommand(rdrTree, {});
        });
};
