import { Text, View } from "@/components/Themed";
import { Stack } from "expo-router";

export default function Page() {
	return (
		<View
			style={{ flex: 1, justifyContent: "center", alignItems: "center" }}
		>
			<Stack.Screen
				name="keys"
				options={{
					title: "keys",
				}}
			/>
			<Text>keys tbd</Text>
		</View>
	);
}
