import { PostCard } from "@/components/card";
import { useEffect, useState } from "react";

// 時間格式轉換函式
function formatTime(isoString) {
    if (!isoString) return "";
    const date = new Date(isoString);
    // 這裡用本地時間格式，也可自訂
    return date.toLocaleString("zh-TW", {
        year: "numeric",
        month: "2-digit",
        day: "2-digit",
        hour: "2-digit",
        minute: "2-digit",
        hour12: false,
    });
}

export default function Feed({ feedData = [] }) {
    const [authorProfiles, setAuthorProfiles] = useState({});

    useEffect(() => {
        const uniqueAuthorIds = [...new Set(feedData.map(post => post.author_id))];
        if (uniqueAuthorIds.length === 0) return;

        Promise.all(
            uniqueAuthorIds.map(async (author_id) => {
                try {
                    const res = await fetch(
                        `${process.env.NEXT_PUBLIC_API_BASE_URL}${process.env.NEXT_PUBLIC_PROFILE_API}${author_id}`,
                        {
                            method: "GET",
                            credentials: 'include',
                        }
                    );

                    if (!res.ok) {
                        throw new Error(`Failed to fetch profile for author_id ${author_id}`);
                    }

                    const profile = await res.json();
                    const avatar_url = profile.avatar_url || null;
                    return [author_id, { ...profile, avatar_url }];
                } catch {
                    return [author_id, null];
                }
            })
        ).then(results => {
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
                    post={{
                        ...post,
                        // 假設 post.created_at 為 ISO 字串
                        created_at: formatTime(post.created_at),
                    }}
                    authorProfile={authorProfiles[post.author_id]}
                />
            ))}
        </div>
    );
}