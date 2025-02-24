import * as FileSystem from "expo-file-system";
import { PersistStorage, StorageValue } from "zustand/middleware";

const allPrefix = `${FileSystem.documentDirectory}gutter`;

export class ZustandStorage<T> implements PersistStorage<T> {
	public dir: string = "";
	public fileSuffix: string = "";

	static getItemPath(
		dir: string,
		filename: string,
		fileSuffix: string,
	): string {
		return `${allPrefix}/${dir}/${filename}${fileSuffix}`;
	}

	static create<T>(dir: string, fileSuffix: string): ZustandStorage<T> {
		const s = new ZustandStorage<T>();
		console.log("initializing new storage", "dir", dir);
		try {
			FileSystem.makeDirectoryAsync(allPrefix).finally(() => {
				FileSystem.makeDirectoryAsync(`${allPrefix}/${dir}`)
					.then(() => {
						console.log("created new dir", dir);
					})
					.catch(async (e) => {
						console.log(
							"creating new dir failed1, listing dirs",
							e,
						);
						FileSystem.readDirectoryAsync(`${allPrefix}/${dir}`)
							.then((v) => {
								console.log("read dir async:", v);
							})
							.catch((e) => {
								console.log("read dir failed:", e);
							});
					});
			});
		} catch (err: any) {
			console.log("creating new dir failed2", err);
		}
		s.dir = dir;
		s.fileSuffix = fileSuffix;
		return s;
	}

	getItem = async (name: string): Promise<StorageValue<T> | null> => {
		const result = await FileSystem.readAsStringAsync(
			`${allPrefix}/${this.dir}/${name}${this.fileSuffix}`,
		);
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
		return await FileSystem.writeAsStringAsync(
			`${allPrefix}/${this.dir}/${name}${this.fileSuffix}`,
			JSON.stringify(value.state, null, 2),
		);
	};

	removeItem = async (name: string): Promise<unknown> => {
		return await FileSystem.deleteAsync(
			`${allPrefix}/${this.dir}/${name}${this.fileSuffix}`,
		);
	};
}
