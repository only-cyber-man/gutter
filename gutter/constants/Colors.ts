const tintColorLight = "#2f95dc";
const tintColorDark = "#fff";

export const ColorValues = {
	Red: "#FF0000",
	Gold: "#ffaa11",
	Amethyst: "#9B5DE5",
	BrilliantRose: "#F15BB5",
	Maize: "#FEE440",
	DeepSkyBlue: "#00BBF9",
	Aquamarine: "#00F5D4",

	LightRed: "#ff5e71",
	LightPink: "#ffd1d6",
	Honey: "#f6d110",
	Butter: "#fff9c7",
	Sky: "#e6f2f4",
	Blue: "#81ceeb",
	Green: "#00ab66",

	DarkBlue: "#00296b",
};

export const Colors = {
	light: {
		text: "#000",
		background: "#fff",
		buttonBackground: ColorValues.Blue,
		buttonText: ColorValues.Sky,
		tint: tintColorLight,
		tabIconDefault: "#ccc",
		tabIconSelected: tintColorLight,
	},
	dark: {
		text: "#fff",
		background: "#000",
		buttonBackground: ColorValues.Blue,
		buttonText: "#000",
		tint: tintColorDark,
		tabIconDefault: "#ccc",
		tabIconSelected: tintColorDark,
	},
};
