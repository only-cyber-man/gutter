import { useUser } from "@/hooks/useUser";
import { Button, ButtonText, View, YStack } from "../Themed";
import { router } from "expo-router";
import { ActivityIndicator, Alert } from "react-native";
import { ColorValues } from "@/constants/Colors";

export const AccountBox = () => {
	const { logout, deleteAccount, token, isLoading } = useUser();

	return (
		<YStack>
			{isLoading && (
				<>
					<ActivityIndicator />
					<View style={{ marginBottom: 25 }} />
				</>
			)}
			<Button
				style={{
					backgroundColor: ColorValues.LightRed,
				}}
				onPress={() => {
					Alert.alert(
						"Are you sure you want to do this?",
						"This operation can't be reverted!",
						[
							{
								text: "Cancel",
								isPreferred: true,
							},
							{
								text: "Yes",
								onPress: async () => {
									await deleteAccount(token);
									router.dismissAll();
								},
							},
						],
					);
				}}
			>
				<ButtonText
					style={{
						color: "white",
					}}
				>
					Delete my account
				</ButtonText>
			</Button>
			<Button
				onPress={() => {
					logout();
					router.dismissAll();
				}}
			>
				<ButtonText>Log out</ButtonText>
			</Button>
		</YStack>
	);
};
