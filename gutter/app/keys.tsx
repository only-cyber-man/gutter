import { UserKeyList } from "@/components/keys/UserKeyList";
import { View } from "@/components/Themed";
import { Stack } from "expo-router";

export default function Page() {
	return (
		<View
			style={{ flex: 1, justifyContent: "center", alignItems: "center" }}
		>
			<Stack.Screen
				name="user keys"
				options={{
					title: "user keys",
				}}
			/>
			<UserKeyList />
		</View>
	);
}
