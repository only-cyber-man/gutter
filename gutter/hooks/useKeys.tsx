import { create } from "zustand";
import { persist } from "zustand/middleware";
import { RSA } from "react-native-rsa-native";
import { ZustandStorage } from "./ZustandStorage";

// Generate RSA key pair in base64 format
export const generateKeyPairRSA = async (): Promise<KeyPair> => {
	return await RSA.generateKeys(2048);
};

// Encrypt a message using the recipient's public key
export const encryptMessageRSA = async (
	message: string,
	recipientPubKey: string,
): Promise<string> => {
	return await RSA.encrypt(message, recipientPubKey);
};
export const encryptLongMessageRSA = async (
	message: string,
	recipientPubKey: string,
): Promise<string> => {
	const splitted = message.match(/.{1,256}/g)!;
	const encrypted: string[] = [];
	for (const s of splitted) {
		encrypted.push(await encryptMessageRSA(s, recipientPubKey));
	}
	return encrypted.join(";");
};

export const decryptMessageRSA = async (
	encrypted: string,
	recipientPrivKey: string,
): Promise<string> => {
	return await RSA.decrypt(encrypted, recipientPrivKey);
};

export const decryptLongMessageRSA = async (
	longEncrypted: string,
	recipientPrivKey: string,
): Promise<string> => {
	const splitted = longEncrypted.split(";");
	let message = "";
	for (const s of splitted) {
		const decryptedChunk = await decryptMessageRSA(s, recipientPrivKey);
		message += decryptedChunk;
	}
	return message;
};

export interface KeyPair {
	public: string;
	private: string;
}

export interface KeysStore {
	userKeys: Record<string, KeyPair>;
	isLoading: boolean;
	chats: Record<string, KeyPair>;
	saveUserPair: (userId: string, newKeyPair: KeyPair) => void;
	addNewChat: (chatId: string, keyPair: KeyPair) => void;
}

export const useKeys = create(
	persist<KeysStore>(
		(set, get) => ({
			chats: {},
			isLoading: false,
			userKeys: {},
			saveUserPair: (userId: string, newKeyPair: KeyPair) => {
				const { userKeys } = get();
				set({
					userKeys: {
						...userKeys,
						[userId]: newKeyPair,
					},
				});
			},
			addNewChat: (chatId, keyPair) => {
				const { chats: _chats } = get();
				set({
					chats: {
						..._chats,
						[chatId]: keyPair,
					},
				});
			},
		}),
		{
			name: "gutter-keys-store",
			storage: ZustandStorage.create("keys", ".key.json"),
		},
	),
);
