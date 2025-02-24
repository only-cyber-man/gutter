import { KeyList } from "@/components/keys/KeyList";
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
			<KeyList />
		</View>
	);
}
