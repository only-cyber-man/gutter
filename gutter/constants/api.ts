export const baseUrl = "https://gutter.cyber-man.pl/api";

export const buildHeaders = (token?: string): { [key: string]: string } => {
	return {
		...(token
			? {
					Authorization: token,
				}
			: {}),
		"Content-Type": "application/json",
	};
};
