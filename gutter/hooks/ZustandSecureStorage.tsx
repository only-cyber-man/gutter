import { PersistStorage, StorageValue } from "zustand/middleware";
import * as SecureStore from "expo-secure-store";

export class ZustandSecureStorage<T> implements PersistStorage<T> {
	getItem = async (name: string): Promise<StorageValue<T> | null> => {
		const result = await SecureStore.getItemAsync(name);
		if (!result) {
			return null;
		}
		return {
			state: JSON.parse(result),
		};
	};
	setItem = async (
		name: string,
		value: StorageValue<T>,
	): Promise<unknown> => {
		return await SecureStore.setItemAsync(
			name,
			JSON.stringify(value.state),
		);
	};
	removeItem = (name: string): Promise<unknown> => {
		return SecureStore.deleteItemAsync(name);
	};
}
