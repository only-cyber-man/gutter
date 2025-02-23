import { useEffect, useState } from "react";
import { Button, ButtonText, Text, XStack, YStack } from "../Themed";
import { ActivityIndicator, FlatList, TouchableOpacity } from "react-native";
import { router } from "expo-router";
import { MaterialCommunityIcons } from "@expo/vector-icons";
import { useUser } from "@/hooks/useUser";
import {
	KeyExchanges,
	FriendStatuses,
	TranslateFriendshipStatus,
	useFriends,
} from "@/hooks/useFriends";
import { ColorValues } from "@/constants/Colors";
import { decryptLongMessageRSA, useKeys } from "@/hooks/useKeys";

export const FriendshipsList = () => {
	const { token, user } = useUser();
	const { getFriendships, isLoading, answerFriend } = useFriends();
	const { userKeys, addNewChat } = useKeys();
	const [friendships, setFriendships] = useState<KeyExchanges[]>([]);

	const refresh = async () => {
		const friendships = await getFriendships(token);
		setFriendships(friendships);
	};

	useEffect(() => {
		refresh();
	}, [token]);

	const renderItem = ({ item }: { item: KeyExchanges; index: number }) => {
		const isSender = user?.id === item.requester;

		return (
			<XStack
				key={item.id}
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
					{isSender
						? item.expand.target?.username
						: item.expand.requester?.username}
				</Text>
				{/* if the user is the requester, he can just wait */}
				<Text
					style={{
						marginRight: 8,
					}}
				>
					{TranslateFriendshipStatus(item.expand.friendship!.status)}
				</Text>
				{/* if the user is the invitee, he can accept the friendship */}
				{!isSender &&
					item.expand.friendship!.status !==
						FriendStatuses.Friends && (
						<TouchableOpacity
							onPress={async () => {
								try {
									const decodedPrivKey =
										await decryptLongMessageRSA(
											item.encryptedPrivateKey,
											userKeys[user!.id]!.private,
										);
									await answerFriend(
										item.expand.friendship!.id,
										true,
										token,
									);
									addNewChat(item.relatedChat, {
										public: item.expand.relatedChat!
											.publicKey,
										private: decodedPrivKey,
									});
									await refresh();
								} catch (err: any) {
									console.log("answer err", err);
								}
							}}
							disabled={isLoading}
						>
							<MaterialCommunityIcons
								name="check-circle-outline"
								size={32}
								color={ColorValues.Green}
							/>
						</TouchableOpacity>
					)}
				<TouchableOpacity
					onPress={async () => {
						await answerFriend(
							item.expand.friendship!.id,
							false,
							token,
						);
						await refresh();
					}}
					disabled={isLoading}
				>
					<MaterialCommunityIcons
						name="close-circle-outline"
						size={32}
						color={ColorValues.Red}
					/>
				</TouchableOpacity>
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
			<XStack
				style={{
					width: "100%",
					justifyContent: "center",
				}}
			>
				<Button
					onPress={() => {
						router.navigate("/invite-friend");
					}}
				>
					<ButtonText>Invite</ButtonText>
				</Button>
				{isLoading && (
					<ActivityIndicator
						style={{
							marginLeft: 8,
						}}
					/>
				)}
			</XStack>
			<FlatList
				style={{
					width: "100%",
				}}
				data={friendships}
				renderItem={renderItem}
				refreshing={isLoading}
				onRefresh={refresh}
				ListEmptyComponent={
					<YStack
						style={{
							marginTop: 16,
						}}
					>
						{isLoading ? (
							<ActivityIndicator />
						) : (
							<Text>
								You have not invited any friends yet! Click the
								button above to invite new friends.
							</Text>
						)}
					</YStack>
				}
			/>
		</YStack>
	);
};
