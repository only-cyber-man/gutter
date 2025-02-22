import * as H from "expo-haptics";

export const Haptics = {
	warn: () => {
		H.notificationAsync(H.NotificationFeedbackType.Warning);
	},
	error: () => {
		H.notificationAsync(H.NotificationFeedbackType.Error);
	},
	soft: () => {
		H.impactAsync(H.ImpactFeedbackStyle.Soft);
	},
	light: () => {
		H.impactAsync(H.ImpactFeedbackStyle.Light);
	},
};
