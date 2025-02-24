import { Text, View, XStack, YStack } from "../Themed";
import { router } from "expo-router";
import {
	GestureResponderEvent,
	StyleProp,
	TextStyle,
	TouchableOpacity,
	ViewStyle,
} from "react-native";
import { ColorValues } from "@/constants/Colors";

export const NavigationBox = () => {
	const Kafelek = ({
		name,
		onPress,
		outerStyle,
		innerStyle,
		textStyle,
	}: {
		name: string;
		onPress?: (event: GestureResponderEvent) => void;
		outerStyle?: StyleProp<ViewStyle>;
		innerStyle?: StyleProp<ViewStyle>;
		textStyle?: StyleProp<TextStyle>;
	}) => {
		return (
			<TouchableOpacity
				onPress={onPress}
				style={[
					{
						flex: 1,
						flexDirection: "column",
						justifyContent: "center",
						alignItems: "center",
						height: "100%",
						paddingVertical: 2.5,
					},
					outerStyle,
				]}
			>
				<View
					style={[
						{
							backgroundColor: "blue",
							justifyContent: "center",
							alignItems: "center",
							flex: 1,
							width: "100%",
							borderRadius: 8,
						},
						innerStyle,
					]}
				>
					<Text
						style={[
							{
								fontSize: 26,
								margin: 8,
								textAlign: "center",
								fontWeight: "bold",
							},
							textStyle,
						]}
					>
						{name}
					</Text>
				</View>
			</TouchableOpacity>
		);
	};

	return (
		<YStack
			style={{
				flex: 1,
				justifyContent: "center",
				alignItems: "flex-start",
			}}
		>
			<XStack style={{ flex: 1 }}>
				<Kafelek
					name="user keys"
					onPress={() => {
						router.navigate("/keys");
					}}
					outerStyle={{
						marginHorizontal: 5,
					}}
					innerStyle={{
						backgroundColor: ColorValues.Honey,
					}}
					textStyle={{
						color: ColorValues.Butter,
					}}
				/>
				<Kafelek
					name="chats"
					onPress={() => {
						router.navigate("/chats");
					}}
					outerStyle={{
						marginRight: 5,
					}}
					innerStyle={{
						backgroundColor: ColorValues.Butter,
					}}
					textStyle={{
						color: ColorValues.Honey,
					}}
				/>
			</XStack>
			<XStack style={{ flex: 1 }}>
				<Kafelek
					name="friendships"
					onPress={() => {
						router.navigate("/friendships");
					}}
					outerStyle={{
						marginHorizontal: 5,
					}}
					innerStyle={{
						backgroundColor: ColorValues.Sky,
					}}
					textStyle={{
						color: ColorValues.Blue,
					}}
				/>
				<Kafelek
					name="account"
					onPress={() => {
						router.navigate("/account");
					}}
					outerStyle={{
						marginRight: 5,
					}}
					innerStyle={{
						backgroundColor: ColorValues.Blue,
					}}
					textStyle={{
						color: ColorValues.Sky,
					}}
				/>
			</XStack>
		</YStack>
	);
};
