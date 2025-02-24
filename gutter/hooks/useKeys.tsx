import * as FileSystem from "expo-file-system";
import * as DocumentPicker from "expo-document-picker";
import { saveDocuments } from "@react-native-documents/picker";
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

export interface UserKeyPair extends KeyPair {
	username: string;
}

export interface KeysStore {
	userKeys: Record<string, UserKeyPair>;
	isLoading: boolean;
	chats: Record<string, KeyPair>;
	saveUserPair: (username: string, newKeyPair: KeyPair) => void;
	deleteUserPair: (username: string) => void;
	downloadUserPair: (username: string) => Promise<void>;
	uploadNewUserPair: () => Promise<UserKeyPair>;
	addNewChat: (chatId: string, keyPair: KeyPair) => void;
}

export const useKeys = create(
	persist<KeysStore>(
		(set, get) => ({
			chats: {},
			isLoading: false,
			userKeys: {},
			saveUserPair: (username: string, newKeyPair: KeyPair) => {
				const { userKeys } = get();
				set({
					userKeys: {
						...userKeys,
						[username]: { username, ...newKeyPair },
					},
				});
			},
			deleteUserPair: (username) => {
				const { userKeys } = get();
				delete userKeys[username];
				set({
					userKeys,
				});
			},
			downloadUserPair: async (username) => {
				const all = get();
				const tmpOut = `${FileSystem.cacheDirectory}${username}.keys.json`;
				console.log({ tmpOut });
				await FileSystem.writeAsStringAsync(
					tmpOut,
					JSON.stringify(all.userKeys[username], null, 2),
				);
				await saveDocuments({
					sourceUris: [tmpOut],
					mimeType: "application/json",
					copy: true,
					fileName: `${username}.keys.json`,
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
			uploadNewUserPair: async () => {
				// const [picked] = await pick({
				// 	mode: "open",
				// });
				const doc = await DocumentPicker.getDocumentAsync({
					type: "application/json",
					copyToCacheDirectory: true,
				});
				if (doc.canceled) {
					return;
				}
				const content = await FileSystem.readAsStringAsync(
					doc.assets[0].uri,
				);
				const kp = JSON.parse(content);
				if (!kp.username || !kp.private || !kp.public) {
					throw new Error(
						"key is not in valid form; must have 'username', 'private' and 'public'",
					);
				}
				const { userKeys } = get();
				set({
					userKeys: {
						...userKeys,
						[kp.username]: kp,
					},
				});
				return kp;
			},
		}),
		{
			name: "gutter-keys-store",
			storage: ZustandStorage.create("keys", ".key.json"),
		},
	),
);
