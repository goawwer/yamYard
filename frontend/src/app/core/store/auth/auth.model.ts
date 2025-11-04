export interface Login {
    email: string;
    password: string;
}

export interface SignUp {
    username: string;
    email: string;
    password: string;
    bio?: string;
}