// logout.js (for use in React components)
const logoutUrl = "http://192.168.2.13:8080/api/v1/auth/logout";
export async function logout(authContext, router) {
    console.log("Logging out...");
    const storedToken = typeof window !== "undefined" ? localStorage.getItem('jwtToken') : null;

    if (!storedToken) {
        if (authContext && typeof authContext.logout === 'function') {
            authContext.logout();
        }
        if (router) router.push('/');
        return;
    }

    try {
        const response = await fetch({logoutUrl}, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${storedToken}`,
            },
        });

        if (!response.ok) {
            let errorData;
            try {
                errorData = await response.json();
            } catch (e) {
                errorData = { message: response.statusText };
            }
            if (authContext && typeof authContext.logout === 'function') {
                authContext.logout();
            }
            if (router) router.push('/');
            return;
        }

        try {
            await response.json();
        } catch (e) {
            // Ignore JSON parse errors for empty response
        }

        if (authContext && typeof authContext.logout === 'function') {
            authContext.logout();
        }
        if (router) router.push('/');

    } catch (error) {
        if (authContext && typeof authContext.logout === 'function') {
            authContext.logout();
        }
        if (router) router.push('/');
    }
}