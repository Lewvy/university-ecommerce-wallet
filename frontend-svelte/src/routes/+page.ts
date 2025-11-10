import { redirect } from '@sveltejs/kit';

const checkAuthentication = async () => {
	const isAuthenticated = false;
	return isAuthenticated;
};

export const load = async () => {
	const isAuthenticated = await checkAuthentication();

	if (isAuthenticated) {
		throw redirect(302, '/dashboard');
	}

	throw redirect(302, '/login');
};
