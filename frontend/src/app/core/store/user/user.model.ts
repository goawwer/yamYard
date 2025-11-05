export interface User {
    id: string;
    email: string;
    username: string;
    hashedPassword: string;
    bio?: string;
    profileStatus?: string;
    image_url?: string;
    is_admin: string;
    createdAt: string;
    updatedAt: string;
}

export interface isUpdatingUser {
    username: string | null | undefined;
    avatar: string | null | undefined;
    profileStatus: string | null | undefined;
    bio: string | null | undefined;
    image_url: string | null | undefined;
}