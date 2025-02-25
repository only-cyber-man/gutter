import { ActivityIndicator, Alert, FlatList } from "react-native";
import { Button, ButtonText, Text, XStack, YStack } from "../Themed";
import { useKeys } from "../../hooks/useKeys";
import { ColorValues } from "../../constants/Colors";

export const UserKeyList = () => {
	const { userKeys, isLoading, deleteUserPair, downloadUserPair } = useKeys();

	const renderItem = ({
		item: username,
		index,
	}: {
		item: string;
		index: number;
	}) => {
		return (
			<XStack
				key={`kp-${index}`}
				style={{
					borderBottomWidth: 1,
					borderBottomColor: "gray",
					width: "100%",
					padding: 24,
				}}
			>
				<Text
					style={{
						fontWeight: "bold",
						flex: 1,
					}}
				>
					{username}
				</Text>
				<Button
					onPress={async () => {
						try {
							await downloadUserPair(username);
						} catch (err) {
							console.log("err", err);
						}
					}}
				>
					<ButtonText>Download</ButtonText>
				</Button>
				<Button
					style={{
						marginLeft: 2,
						backgroundColor: ColorValues.LightRed,
					}}
					onPress={() => {
						Alert.alert(
							"Are you sure?",
							"this operation is irreversible. Once you loose the key without backup, you won't be able to access that account.",
							[
								{
									text: "Cancel",
									isPreferred: true,
									style: "cancel",
								},
								{
									text: "I'm sure, delete the key",
									onPress: () => {
										deleteUserPair(username);
									},
									style: "default",
								},
							],
						);
					}}
				>
					<ButtonText
						style={{
							color: "black",
						}}
					>
						Delete
					</ButtonText>
				</Button>
			</XStack>
		);
	};

	return (
		<YStack
			style={{
				width: "100%",
				flex: 1,
				marginTop: 8,
			}}
		>
			<FlatList
				style={{
					width: "100%",
				}}
				data={Object.keys(userKeys)}
				renderItem={renderItem}
				ListEmptyComponent={
					<YStack
						style={{
							marginTop: 16,
						}}
					>
						{isLoading ? (
							<ActivityIndicator />
						) : (
							<Text>Please log out to see the keys</Text>
						)}
					</YStack>
				}
			/>
		</YStack>
	);
};
