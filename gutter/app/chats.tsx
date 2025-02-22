import { Text, View } from "@/components/Themed";
import { Stack } from "expo-router";

export default function Page() {
	return (
		<View
			style={{ flex: 1, justifyContent: "center", alignItems: "center" }}
		>
			<Stack.Screen
				name="chats"
				options={{
					title: "chats",
				}}
			/>
			<Text>chats tbd</Text>
		</View>
	);
}
