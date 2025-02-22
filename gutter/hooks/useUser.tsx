import { create } from "zustand";
import { persist } from "zustand/middleware";
import { Alert } from "react-native";
import { getNotificationsToken } from "@/api/notifications";
import { baseUrl, buildHeaders } from "@/constants/api";
import { ZustandStorage } from "./ZustandStorage";

export interface User {
	id: string;
	username: string;
	publicKey: string;
}

export interface UserStore {
	token: string;
	isLoading: boolean;
	user: User | null;

	isLoggedIn: () => boolean;
	login: (username: string, password: string) => Promise<void>;
	register: (
		username: string,
		password: string,
		publicKey: string,
	) => Promise<void>;
	logout: () => void;
	deleteAccount: (token: string) => Promise<void>;
}

export const useUser = create(
	persist<UserStore>(
		(set, get) => ({
			token: "",
			user: null,
			isLoading: false,

			isLoggedIn: () => {
				const { token, user } = get();
				return token !== "" && user !== null;
			},
			login: async (username, password) => {
				try {
					set({ isLoading: true });
					let pushToken: string | undefined;
					try {
						pushToken = await getNotificationsToken();
					} catch (err) {
						console.log("failed to get expo token", err);
					}
					const response = await fetch(`${baseUrl}/auth/login`, {
						method: "POST",
						body: JSON.stringify({
							username,
							password,
							pushToken,
						}),
						headers: buildHeaders(),
					});
					const { data, message, success } = await response.json();
					if (!success) {
						throw new Error(message);
					}
					const { token, user } = data;
					set({
						token,
						user,
					});
				} catch (err: any) {
					Alert.alert("Error occurred", err.message);
				} finally {
					set({ isLoading: false });
				}
			},
			register: async (username, password, publicKey) => {
				try {
					set({ isLoading: true });
					let pushToken: string | undefined;
					try {
						pushToken = await getNotificationsToken();
						console.log(pushToken);
					} catch (err) {
						console.log("failed to get expo token", err);
					}
					const response = await fetch(`${baseUrl}/auth/register`, {
						method: "POST",
						body: JSON.stringify({
							username,
							password,
							pushToken,
							publicKey,
						}),
						headers: buildHeaders(),
					});
					const { data, message, success } = await response.json();
					if (!success) {
						throw new Error(message);
					}
					const { token, user } = data;
					set({
						token,
						user,
					});
				} catch (err: any) {
					Alert.alert("Error occurred", err.message);
				} finally {
					set({ isLoading: false });
				}
			},
			logout: () => {
				set({
					token: "",
					user: null,
				});
			},
			deleteAccount: async (token) => {
				try {
					set({ isLoading: true });
					const response = await fetch(`${baseUrl}/auth/account`, {
						method: "DELETE",
						headers: buildHeaders(token),
					});
					const { message, success } = await response.json();
					if (!success) {
						throw new Error(message);
					}
					get().logout();
				} catch (err: any) {
					Alert.alert("Error occurred", err.message);
				} finally {
					set({ isLoading: false });
				}
			},
		}),
		{
			name: "gutter-user-store",
			storage: new ZustandStorage(),
		},
	),
);
