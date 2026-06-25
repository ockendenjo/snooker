export interface Environment {
    production: boolean;
    startDate?: Date;
    endDate?: Date;
    cognito: {
        domain: string;
        clientId: string;
    };
}
