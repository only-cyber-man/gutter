import { FriendshipsList } from "@/components/friendships/FriendshipsList";
import { View } from "@/components/Themed";
import { Stack } from "expo-router";

export default function Page() {
	return (
		<View
			style={{ flex: 1, justifyContent: "center", alignItems: "center" }}
		>
			<Stack.Screen
				name="friendships"
				options={{
					title: "friendships",
				}}
			/>
			<FriendshipsList />
		</View>
	);
}
