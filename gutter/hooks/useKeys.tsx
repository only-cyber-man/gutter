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

export const decryptMessageRSA = async (
	encrypted: string,
	recipientPrivKey: string,
): Promise<string> => {
	return await RSA.decrypt(encrypted, recipientPrivKey);
};

export interface KeyPair {
	public: string;
	private: string;
}

export interface KeysStore {
	user: KeyPair | null;
	isLoading: boolean;
	chats: KeyPair[];
	createNewUserKeypair: () => Promise<KeyPair>;
}

export const useKeys = create(
	persist<KeysStore>(
		(set, get) => ({
			chats: [],
			isLoading: false,
			user: null,
			createNewUserKeypair: async () => {
				const newPair = await generateKeyPairRSA();
				set({
					user: newPair,
				});
				return newPair;
			},
		}),
		{
			name: "gutter-keys-store",
			storage: new ZustandStorage(),
		},
	),
);
