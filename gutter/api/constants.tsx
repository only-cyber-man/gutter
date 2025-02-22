import Constants from "expo-constants";

export const getProjectId = (): string => {
	return (
		Constants?.expoConfig?.extra?.eas?.projectId ??
		Constants?.easConfig?.projectId
	);
};
