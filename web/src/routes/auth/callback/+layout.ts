// Disable SSR and parent layout load functions for callback page
// This prevents auth redirect while processing the OAuth callback
export const ssr = false;

export const load = () => {
	return {};
};
