// logout.js (for use in React components)

export async function logout(authContext, router) {
    console.log("Logging out...");

    try {
        const response = await fetch(
            `${process.env.NEXT_PUBLIC_API_BASE_URL}${process.env.NEXT_PUBLIC_LOGOUT_API}`, {
            method: 'POST',
            credentials: 'include', // 讓瀏覽器帶上 cookie
        });

        // 不論 response 狀態如何，都清除前端狀態並導回首頁
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