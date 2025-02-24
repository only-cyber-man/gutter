import { generateKeyPairRSA, useKeys } from "@/hooks/useKeys";
import { Button, ButtonText, XStack, YStack, TextInput, Text } from "../Themed";
import { useState } from "react";
import { ActivityIndicator, Alert } from "react-native";
import { useUser } from "@/hooks/useUser";

export const AuthBox = () => {
	const { saveUserPair, userKeys, uploadNewUserPair } = useKeys();
	const { login, register, isLoading } = useUser();
	const [username, setUsername] = useState("");

	return (
		<YStack
			style={{
				padding: 24,
			}}
		>
			<TextInput
				placeholder="username..."
				textContentType="username"
				onChangeText={setUsername}
			/>
			<XStack style={{ width: "100%", margin: 24 }}>
				<Button
					disabled={isLoading}
					onPress={async () => {
						try {
							if (!userKeys[username]) {
								Alert.alert(
									"You don't have keys to that account",
								);
								return;
							}
							await login(username, userKeys[username].private);
						} catch (err) {
							console.log("err", err);
						}
					}}
				>
					<ButtonText>login</ButtonText>
				</Button>
				<Button
					disabled={isLoading}
					onPress={async () => {
						try {
							const keypair = await generateKeyPairRSA();
							const { user } = await register(
								username,
								keypair.public,
								keypair.private,
							);
							saveUserPair(user.username, keypair);
						} catch (err) {
							console.log("err", err);
						}
					}}
				>
					<ButtonText>register</ButtonText>
				</Button>
			</XStack>
			<XStack>
				<Text style={{ flex: 1 }}>
					do you have a key from the another device?
				</Text>
				<Button
					onPress={async () => {
						try {
							const kp = await uploadNewUserPair();
							await login(kp.username, kp.private);
						} catch (err) {
							console.log("err", err);
						}
					}}
				>
					<ButtonText>import</ButtonText>
				</Button>
			</XStack>
			{isLoading && (
				<XStack>
					<ActivityIndicator />
				</XStack>
			)}
		</YStack>
	);
};
