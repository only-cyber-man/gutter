/**
 * Learn more about Light and Dark modes:
 * https://docs.expo.io/guides/color-schemes/
 */

import {
	Text as DefaultText,
	View as DefaultView,
	TouchableOpacity as DefaultTouchableOpacity,
	TextInput as DefaultTextInput,
	TouchableOpacityProps as DefaultTouchableOpacityProps,
} from "react-native";

import { Colors } from "@/constants/Colors";
import { useColorScheme } from "./useColorScheme";
import { Haptics } from "@/api/haptics";
import { MaterialCommunityIcons as DefaultIcon } from "@expo/vector-icons";
import { IconProps as DefaultIconProps } from "@expo/vector-icons/build/createIconSet";

type ThemeProps = {
	lightColor?: string;
	darkColor?: string;
};

export type IconProps = ThemeProps & DefaultIconProps<any>;
export type TextProps = ThemeProps & DefaultText["props"];
export type TouchableOpacityProps = ThemeProps & DefaultTouchableOpacityProps;
export type ViewProps = ThemeProps & DefaultView["props"];
export type TextInputProps = ThemeProps & DefaultTextInput["props"];

export function useThemeColor(
	props: { light?: string; dark?: string },
	colorName: keyof typeof Colors.light & keyof typeof Colors.dark,
) {
	const theme = useColorScheme() ?? "light";
	const colorFromProps = props[theme];

	if (colorFromProps) {
		return colorFromProps;
	} else {
		return Colors[theme][colorName];
	}
}

export function Text(props: TextProps) {
	const { style, lightColor, darkColor, ...otherProps } = props;
	const color = useThemeColor({ light: lightColor, dark: darkColor }, "text");

	return <DefaultText style={[{ color }, style]} {...otherProps} />;
}

export function Icon(props: IconProps) {
	const { lightColor, darkColor, ...otherProps } = props;
	const color = useThemeColor({ light: lightColor, dark: darkColor }, "text");

	return <DefaultIcon color={color} {...otherProps} />;
}

export function View(props: ViewProps) {
	const { style, lightColor, darkColor, ...otherProps } = props;
	const backgroundColor = useThemeColor(
		{ light: lightColor, dark: darkColor },
		"background",
	);

	return <DefaultView style={[{ backgroundColor }, style]} {...otherProps} />;
}

export function Button(props: TouchableOpacityProps) {
	const { style, lightColor, darkColor, ...otherProps } = props;
	const buttonBackground = useThemeColor(
		{ light: lightColor, dark: darkColor },
		"buttonBackground",
	);
	return (
		<DefaultTouchableOpacity
			onPress={(e) => {
				Haptics.light();
				if (props.onPress) {
					props.onPress(e);
				}
			}}
			style={[
				{
					padding: 8,
					backgroundColor: buttonBackground,
					borderRadius: 3,
					marginBottom: 4,
				},
				style,
			]}
			{...otherProps}
		/>
	);
}

export function ButtonText(props: TextProps) {
	const { style, lightColor, darkColor, ...otherProps } = props;
	const buttonText = useThemeColor(
		{ light: lightColor, dark: darkColor },
		"buttonText",
	);

	return (
		<DefaultText
			style={[
				{
					color: buttonText,
					marginHorizontal: 8,
					fontWeight: "bold",
					fontSize: 16,
				},
				style,
			]}
			{...otherProps}
		/>
	);
}

export function TextInput(props: TextInputProps) {
	const { style, lightColor, darkColor, ...otherProps } = props;
	const text = useThemeColor({ light: lightColor, dark: darkColor }, "text");

	return (
		<DefaultTextInput
			style={[
				{
					padding: 8,
					borderBottomWidth: 1,
					borderBottomColor: text,
					width: "100%",
					fontSize: 18,
					fontFamily: "SpaceMono",
					color: text,
				},
				style,
			]}
			autoCapitalize="none"
			autoCorrect={false}
			{...otherProps}
		/>
	);
}

export function XStack(props: ViewProps) {
	const { style, ...otherProps } = props;
	return (
		<DefaultView
			style={[
				{
					display: "flex",
					flexDirection: "row",
					justifyContent: "space-between",
					alignItems: "center",
				},
				style,
			]}
			{...otherProps}
		/>
	);
}

export function YStack(props: ViewProps) {
	const { style, ...otherProps } = props;
	return (
		<DefaultView
			style={[
				{
					display: "flex",
					flexDirection: "column",
					justifyContent: "center",
					alignItems: "center",
				},
				style,
			]}
			{...otherProps}
		/>
	);
}
