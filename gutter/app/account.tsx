import { AccountBox } from "@/components/accountscreen/AccountBox";
import { View } from "@/components/Themed";
import { Stack } from "expo-router";

export default function Page() {
	return (
		<View style={{ flex: 1, justifyContent: "center" }}>
			<Stack.Screen
				name="account"
				options={{
					title: "account",
				}}
			/>
			<AccountBox />
		</View>
	);
}
