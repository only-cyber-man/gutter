import { baseUrl, buildHeaders } from "@/constants/api";
import { Alert } from "react-native";
import { create } from "zustand";

export interface User {
	id: string;
	username: string;
	publicKey: string;
}

export interface Chat {
	id: string;
	creator: string;
	participants: string[];
	publicKey: string;
	updated: string;
	created: string;

	expand: {
		creator: User | null;
		participants: User[] | null;
	};
}

export enum FriendStatuses {
	RequestSent = "request sent",
	Friends = "friends",
}

export type FriendStatus = FriendStatuses.Friends | FriendStatuses.RequestSent;

export const TranslateFriendshipStatus = (status: FriendStatus): string => {
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
	id: string;
	requester: string;
	invitee: string;
	status: FriendStatus;
	created: string;
	updated: string;

	expand: {
		requester: User | null;
		invitee: User | null;
	};
}

export interface KeyExchanges {
	id: string;
	relatedChat: string;
	requester: string;
	target: string;
	friendship: string;
	encryptedPrivateKey: string;
	updated: string;
	created: string;

	expand: {
		relatedChat: Chat | null;
		requester: User | null;
		target: User | null;
		friendship: Friendship | null;
	};
}

export interface FriendsListStore {
	isLoading: boolean;
	getFriendships: (token?: string) => Promise<KeyExchanges[]>;
	getUserByUsername: (username: string) => Promise<User>;
	inviteFriend: (
		username: string,
		encryptedPrivateKey: string,
		chatPublicKey: string,
		token?: string,
	) => Promise<Chat>;
	answerFriend: (
		friendshipId: string,
		accept: boolean,
		token?: string,
	) => Promise<void>;
}

export const useFriends = create<FriendsListStore>((set, get) => ({
	isLoading: false,
	getUserByUsername: async (username: string): Promise<User> => {
		try {
			set({ isLoading: true });
			const params = new URLSearchParams();
			params.set("username", username);
			const url = `${baseUrl}/find-user?${params.toString()}`;
			const response = await fetch(url, {
				method: "GET",
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
	inviteFriend: async (
		username: string,
		encryptedPrivateKey: string,
		chatPublicKey: string,
		token?: string,
	) => {
		try {
			set({ isLoading: true });
			const response = await fetch(`${baseUrl}/friendships/invite`, {
				method: "POST",
				headers: buildHeaders(token),
				body: JSON.stringify({
					username,
					encryptedPrivateKey,
					chatPublicKey,
				}),
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
