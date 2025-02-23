import { generateKeyPairRSA, useKeys } from "@/hooks/useKeys";
import { Button, ButtonText, XStack, YStack, TextInput } from "../Themed";
import { useState } from "react";
import { ActivityIndicator } from "react-native";
import { useUser } from "@/hooks/useUser";

export const AuthBox = () => {
	const { saveUserPair } = useKeys();
	const { login, register, isLoading } = useUser();
	const [username, setUsername] = useState("");
	const [password, setPassword] = useState("");

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
			<TextInput
				placeholder="password..."
				secureTextEntry
				textContentType="password"
				onChangeText={setPassword}
			/>
			<XStack style={{ width: "100%", margin: 24 }}>
				<Button
					disabled={isLoading}
					onPress={() => {
						login(username, password);
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
								password,
								keypair.public,
							);
							saveUserPair(user.id, keypair);
						} catch (err) {
							console.log("err", err);
						}
					}}
				>
					<ButtonText>register</ButtonText>
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
