import { useFriends } from "@/hooks/useFriends";
import { useState } from "react";
import { Button, ButtonText, TextInput, XStack, YStack } from "../Themed";
import { router } from "expo-router";
import { ActivityIndicator, View } from "react-native";
import { useUser } from "@/hooks/useUser";
import {
	encryptLongMessageRSA,
	generateKeyPairRSA,
	useKeys,
} from "@/hooks/useKeys";

export const InviteFriendForm = ({
	defaultUsername,
}: {
	defaultUsername?: string;
}) => {
	const { token } = useUser();
	const { inviteFriend, getUserByUsername, isLoading } = useFriends();
	const { addNewChat } = useKeys();
	const [username, setUsername] = useState(defaultUsername ?? "");
	return (
		<YStack
			style={{
				padding: 24,
				justifyContent: "flex-start",
			}}
		>
			<TextInput
				placeholder="Friend's username"
				textContentType="username"
				onChangeText={setUsername}
			/>
			<XStack style={{ width: "100%", margin: 24 }}>
				<View style={{ flex: 1 }} />
				<Button
					disabled={isLoading}
					onPress={async () => {
						try {
							// console.log(1);
							const invitee = await getUserByUsername(username);
							// console.log(2);
							const chatKeyPair = await generateKeyPairRSA();
							// console.log(3);
							const encryptedPrivateKey =
								await encryptLongMessageRSA(
									chatKeyPair.private,
									invitee.publicKey,
								);
							// console.log(4);
							const createdChat = await inviteFriend(
								username,
								encryptedPrivateKey,
								chatKeyPair.public,
								token,
							);
							// console.log(5);
							addNewChat(createdChat.id, chatKeyPair);
							// console.log(6);
							router.dismissAll();
							router.navigate("/friendships");
						} catch (err: any) {
							console.log("err", err);
						}
					}}
				>
					<ButtonText>Invite</ButtonText>
				</Button>
			</XStack>
			{isLoading && <ActivityIndicator />}
		</YStack>
	);
};
