import { PostCard } from "@/components/card";
import { useEffect, useState } from "react";

export default function Feed({ feedData = [] }) {
    const [authorProfiles, setAuthorProfiles] = useState({});

    useEffect(() => {

        // 取得所有唯一的 author_id
        const uniqueAuthorIds = [...new Set(feedData.map(post => post.author_id))];
        if (uniqueAuthorIds.length === 0) return;

        // 批次 fetch 所有作者資料
        Promise.all(
            uniqueAuthorIds.map(async (author_id) => {
                try {
                    const res = await fetch(
                        `${process.env.NEXT_PUBLIC_API_BASE_URL}${process.env.NEXT_PUBLIC_PROFILE_API}${author_id}`
                    , {
                        method: "GET",
                        credentials: 'include', 
                    });

                    if (!res.ok){
                        throw new Error(`Failed to fetch profile for author_id ${author_id}`);
                    }

                    const profile = await res.json();
                    // 轉換 avatar_access_key 為 viewableUrl
                    const avatar_url = profile.avatar_url || null;
                    return [author_id, { ...profile, avatar_url }];
                } catch {
                    return [author_id, null];
                }
            })
        ).then(results => {
            // 轉成 { authorId: profileObj }
            const profilesObj = Object.fromEntries(results);
            setAuthorProfiles(profilesObj);
        });
    }, [feedData]);

    return (feedData.length === 0 ?
        <>
            <div className="text-center text-gray-500 p-6">
                <p className="text-lg">No posts available.</p>
            </div>
        </>
        :
        <div>
            {feedData.map((post) => (
                <PostCard
                    key={post.post_id || post.id}
                    post={post}
                    authorProfile={authorProfiles[post.author_id]}
                />
            ))}
        </div>
    );
}