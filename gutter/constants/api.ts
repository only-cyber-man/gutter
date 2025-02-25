// export const baseUrl = "https://gutter.cyber-man.pl/api";
export const baseUrl = "http://192.168.1.193:7005/api";

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

export const checkApi = async () => {
	const response = await fetch(baseUrl);
	if (response.status !== 200) {
		throw new Error(
			"Something went wrong when connecting to the API. Please try again later or contact the administrator.",
		);
	}
};
