import { PersistStorage, StorageValue } from "zustand/middleware";
import AsyncStorage from "@react-native-async-storage/async-storage";

export class ZustandStorage<T> implements PersistStorage<T> {
	getItem = async (name: string): Promise<StorageValue<T> | null> => {
		const result = await AsyncStorage.getItem(name);
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
		return await AsyncStorage.setItem(name, JSON.stringify(value.state));
	};
	removeItem = (name: string): Promise<unknown> => {
		return AsyncStorage.removeItem(name);
	};
}
