export interface User {
    id: string;
    email: string;
    username: string;
    hashedPassword: string;
    bio?: string;
    profileStatus?: string;
    isAdmin: string;
    createdAt: string;
    updatedAt: string;
}