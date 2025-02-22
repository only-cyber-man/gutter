import { View } from "@/components/Themed";
import { Stack, useLocalSearchParams } from "expo-router";
import { InviteFriendForm } from "@/components/friendships/InviteFriendForm";

export default function Page() {
	const { username } = useLocalSearchParams<{
		username: string;
	}>();
	// console.log("using default username", username);

	return (
		<View
			style={{
				flex: 1,
			}}
		>
			<Stack.Screen
				name="invite-friend"
				options={{
					title: "new friend",
				}}
			/>
			<InviteFriendForm defaultUsername={username} />
		</View>
	);
}
