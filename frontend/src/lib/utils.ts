import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

export function userInitial(user: User): string {
	if(user.email.length == 0) {
		return "";
	}

	const initial = user.email.charAt(0).toUpperCase()
	return initial;
}