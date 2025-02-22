import { AuthBox } from "@/components/authentication/AuthBox";
import { NavigationBox } from "@/components/mainscreen/NavigationBox";
import { View } from "@/components/Themed";
import { useUser } from "@/hooks/useUser";
import { Stack, useSegments } from "expo-router";

export default function Page() {
	const { isLoggedIn, user } = useUser();
	const routeSegments = useSegments();

	const getTitle = () => {
		if (!isLoggedIn()) {
			return "register or log in";
		}
		if (routeSegments.length > 0) {
			return "main screen";
		}
		return `welcome ${user?.username}!`;
	};

	return (
		<View style={{ flex: 1 }}>
			<Stack.Screen
				options={{
					title: getTitle(),
				}}
			/>
			{isLoggedIn() ? <NavigationBox /> : <AuthBox />}
		</View>
	);
}
