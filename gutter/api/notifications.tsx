import * as Notifications from "expo-notifications";
import { getProjectId } from "./constants";

export const getNotificationsToken = async (): Promise<string> => {
	try {
		const { status: existingStatus } =
			await Notifications.getPermissionsAsync();
		let finalStatus = existingStatus;
		if (existingStatus !== "granted") {
			const { status } = await Notifications.requestPermissionsAsync();
			finalStatus = status;
		}
		console.log(1, { finalStatus });
		const projectId = getProjectId();
		if (!projectId) {
			throw new Error("projectId is not defined");
		}
		// console.log("using project id", projectId);
		const getTokenResponse = await Notifications.getExpoPushTokenAsync({
			projectId,
		});
		const expoToken = getTokenResponse.data;
		return expoToken;
	} catch (err) {
		console.log("getNotificationsToken failed", err);
		throw err;
	}
};
