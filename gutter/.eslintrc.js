// https://docs.expo.dev/guides/using-eslint/
module.exports = {
	extends: ["expo", "prettier"],
	plugins: ["prettier"],
	ignorePatterns: ["/dist/*", "/ios/*", "expo-env.d.ts"],
	rules: {
		"react-hooks/exhaustive-deps": "off",
		"react-hooks/rules-of-hooks": "off",
		"prettier/prettier": [
			"error",
			{
				endOfLine: "auto",
			},
		],
	},
};
