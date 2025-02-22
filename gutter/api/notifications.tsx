import * as Notifications from "expo-notifications";
import { getProjectId } from "./constants";

export const getNotificationsToken = async (): Promise<string> => {
	const { status: existingStatus } =
		await Notifications.getPermissionsAsync();
	let finalStatus = existingStatus;
	if (existingStatus !== "granted") {
		const { status } = await Notifications.requestPermissionsAsync();
		finalStatus = status;
	}
	console.log({ finalStatus });
	const projectId = getProjectId();
	if (!projectId) {
		throw new Error("projectId is not defined");
	}
	const getTokenResponse = await Notifications.getExpoPushTokenAsync({
		projectId,
	});
	// console.log(getTokenResponse);
	const expoToken = getTokenResponse.data;
	return expoToken;
};
