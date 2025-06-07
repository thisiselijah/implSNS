
export async function GetViewableUrl( avatar_access_key ) {
    if (avatar_access_key) {
        // fetch avatar image URL
        const viewUrlResponse = await fetch(process.env.NEXT_PUBLIC_S3_AVATAR_URL+'?key='+`${encodeURIComponent(avatar_access_key)}`);

        if (!viewUrlResponse.ok) {
            throw new Error('無法獲取讀取連結');
        }

        let res = await viewUrlResponse.json();
        const viewableUrl = res.viewableUrl;
        return viewableUrl;
        
    }
    return null;
}