import { useFriends } from "@/hooks/useFriends";
import { useState } from "react";
import { Button, ButtonText, TextInput, XStack, YStack } from "../Themed";
import { router } from "expo-router";
import { ActivityIndicator, View } from "react-native";
import { useUser } from "@/hooks/useUser";

export const InviteFriendForm = ({
	defaultUsername,
}: {
	defaultUsername?: string;
}) => {
	const { token } = useUser();
	const { inviteFriend, isLoading } = useFriends();
	const [username, setUsername] = useState(defaultUsername ?? "");
	return (
		<YStack
			style={{
				padding: 24,
				justifyContent: "flex-start",
			}}
		>
			<TextInput
				placeholder="Nazwa znajomego"
				textContentType="username"
				onChangeText={setUsername}
			/>
			<XStack style={{ width: "100%", margin: 24 }}>
				<View style={{ flex: 1 }} />
				<Button
					disabled={isLoading}
					onPress={async () => {
						await inviteFriend(username, token);
						router.dismissAll();
						router.navigate("/friendships");
					}}
				>
					<ButtonText>Zaproś</ButtonText>
				</Button>
			</XStack>
			{isLoading && <ActivityIndicator />}
		</YStack>
	);
};
