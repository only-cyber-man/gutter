import { create } from "zustand";
import { persist } from "zustand/middleware";
import { ZustandSecureStorage } from "./ZustandSecureStorage";
import { generateKeyPairSync, publicEncrypt, privateDecrypt } from "crypto";
import { Buffer } from "buffer";
// Generate RSA key pair in base64 format
export const generateKeyPairRSA = (): KeyPair => {
	const { publicKey, privateKey } = generateKeyPairSync("rsa", {
		modulusLength: 2048,
		publicKeyEncoding: { type: "spki", format: "pem" },
		privateKeyEncoding: { type: "pkcs1", format: "pem" },
	});

	return {
		public: Buffer.from(publicKey).toString("base64"),
		private: Buffer.from(privateKey).toString("base64"),
	};
};

// Encrypt a message using the recipient's public key
export const encryptMessageRSA = (message: string, recipientPubKey: string) => {
	const buffer = Buffer.from(message, "utf-8");
	const publicKeyBuffer = Buffer.from(recipientPubKey, "base64").toString(
		"utf-8",
	);
	const encrypted = publicEncrypt(publicKeyBuffer, buffer);
	return encrypted.toString("base64");
};

// Decrypt a message using the recipient's private key
export const decryptMessageRSA = (
	encrypted: string,
	recipientPrivKey: string,
) => {
	const buffer = Buffer.from(encrypted, "base64");
	const privateKeyBuffer = Buffer.from(recipientPrivKey, "base64").toString(
		"utf-8",
	);
	const decrypted = privateDecrypt(privateKeyBuffer, buffer);
	return decrypted.toString("utf-8");
};

export interface KeyPair {
	public: string;
	private: string;
}

export interface KeysStore {
	user: KeyPair;
	isLoading: boolean;
	chats: KeyPair[];
}

export const useKeys = create(
	persist<KeysStore>(
		(set, get) => ({
			chats: [],
			isLoading: false,
			user: generateKeyPairRSA(),
		}),
		{
			name: "gutter-keys-store",
			storage: new ZustandSecureStorage(),
		},
	),
);
