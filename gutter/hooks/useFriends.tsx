import { baseUrl, buildHeaders } from "@/constants/api";
import { Alert } from "react-native";
import { create } from "zustand";

export interface Friend {
	id: string;
	username: string;
}

export interface FriendWithCollections {
	user: Friend;
}

export enum FriendStatuses {
	RequestSent = "request sent",
	Friends = "friends",
}

export const TranslateFriendshipStatus = (status: string): string => {
	switch (status) {
		case FriendStatuses.RequestSent: {
			return "Invite sent";
		}
		case FriendStatuses.Friends: {
			return "Friends";
		}
		default:
			return status;
	}
};

export interface Friendship {
	friendshipId: string;
	invitee: Friend;
	requester: Friend;
	status: string;
}

export interface FriendsListStore {
	isLoading: boolean;
	getFriendships: (token?: string) => Promise<Friendship[]>;
	inviteFriend: (username: string, token?: string) => Promise<void>;
	answerFriend: (
		friendshipId: string,
		accept: boolean,
		token?: string,
	) => Promise<void>;
}

export const useFriends = create<FriendsListStore>((set, get) => ({
	isLoading: false,
	inviteFriend: async (username: string, token?: string) => {
		try {
			set({ isLoading: true });
			const response = await fetch(`${baseUrl}/friendships/invite`, {
				method: "POST",
				headers: buildHeaders(token),
				body: JSON.stringify({
					username,
				}),
			});
			const { message, success } = await response.json();
			if (!success) {
				throw new Error(message);
			}
			set({ isLoading: false });
		} catch (err: any) {
			set({ isLoading: false });
			Alert.alert("Error occurred", err.message);
			throw err;
		}
	},
	answerFriend: async (
		friendshipId: string,
		accept: boolean,
		token?: string,
	) => {
		try {
			set({ isLoading: true });
			const response = await fetch(`${baseUrl}/friendships/answer`, {
				method: "POST",
				headers: buildHeaders(token),
				body: JSON.stringify({
					friendshipId,
					accept,
				}),
			});
			const { message, success } = await response.json();
			if (!success) {
				throw new Error(message);
			}
			set({ isLoading: false });
		} catch (err: any) {
			set({ isLoading: false });
			Alert.alert("Error occurred", err.message);
			throw err;
		}
	},
	getFriendships: async (token) => {
		try {
			set({ isLoading: true });
			const response = await fetch(`${baseUrl}/friendships`, {
				method: "GET",
				headers: buildHeaders(token),
			});
			const { data, message, success } = await response.json();
			if (!success) {
				throw new Error(message);
			}
			set({ isLoading: false });
			return data;
		} catch (err: any) {
			set({ isLoading: false });
			Alert.alert("Error occurred", err.message);
			throw err;
		}
	},
}));
