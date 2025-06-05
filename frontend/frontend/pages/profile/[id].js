import Layout from "@/components/layout";
import { ProfileNavbar } from "@/components/navbar";
import { useState, useEffect, use } from "react";
import Feed from "@/components/feed";

export async function getServerSideProps(context) {
    const { req, params } = context;
    const profileIdFromUrl = params.id; // 從 URL 獲取的動態 ID
    let loggedInUserId = null; // 當前登入用戶的 ID

    if (req && req.cookies) {
        loggedInUserId = req.cookies.user_id || null;
    }

    console.log(`getServerSideProps: Viewing profile for ${profileIdFromUrl}, Logged in as: ${loggedInUserId}`);

    // 接下來，您可以根據業務邏輯處理：
    // 1. 驗證登入用戶是否有權限查看此 profileIdFromUrl 的資料
    // 2. 從資料庫或 API 根據 profileIdFromUrl 獲取個人資料
    //    const profileData = await fetchUserProfile(profileIdFromUrl);
    // 3. 如果找不到資料，可以返回 notFound: true
    //    if (!profileData) {
    //        return { notFound: true };
    //    }

    return {
        props: {
            profileId: profileIdFromUrl,       // 被查看的個人檔案 ID
            currentUserId: loggedInUserId,     // 當前登入用戶的 ID (可用於判斷是否為本人)
            // profileData,                  // 獲取到的個人檔案資料
        },
    };
}

export default function ProfilePage({ profileId, currentUserId /*, profileData */ }) {
    const isOwnProfile = profileId === currentUserId && currentUserId !== null;
    const [editOnClick, setEditOnClick] = useState(false);
    useEffect(() => {
        
    }, []);
    const initialBio = "Hello World"; // 假設這是從後端獲取的初始 bio
    const [bio, setBio] = useState(initialBio); // 假設 bio 初始為空

    useEffect(() => {
        // 這裡可以 fetch 當前 bio
        // (可選) fetch(`/api/profile/${profileId}/bio`).then(...)
    }, [profileId]);

    const handleEditOrDone = async () => {
        if (editOnClick) {
            // 按下 Done，請求後端上傳 bio
            try {
                const res = await fetch(`/api/profile/${profileId}/bio`, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify({ bio }),
                });
                if (!res.ok) {
                    alert("上傳 bio 失敗，請稍後再試");
                    setBio(initialBio); // 如果上傳失敗，恢復初始 bio
                }
            } catch (err) {
                console.error(err);
            }

        }
        setEditOnClick((prev) => !prev);
        
    };
    


    return (
        <Layout pageTitle={isOwnProfile ? "Profile" : `Profile-${profileId}`}>
            <ProfileNavbar />
            <div className="grid grid-cols-1 lg:grid-cols-12 gap-4 px-4 min-h-screen">
                <main className="flex flex-col col-start-4 col-end-10 ">
                    <div className="bg-white p-0.5 rounded-lg shadow">
                        <div className="flex flex-col p-4 border-b border-gray-200">
                            <div className="flex items-center gap-4 p-4">
                                <img
                                    alt="Profile Avatar"
                                    src="https://images.unsplash.com/photo-1491528323818-fdd1faba62cc?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
                                    className="inline-block size-10 rounded-full ring-2 ring-offset-2 ring-[#B6B09F]"
                                />
                                <div>
                                    <h1 className="text-xl font-semibold text-black">
                                        Default name
                                    </h1>
                                    <div className="flex items-center gap-2 text-gray-600">
                                        <span className="text-sm text-gray-600">
                                            #{profileId}
                                        </span>
                                        <span className="text-sm text-gray-600">
                                            Joined on 2023-10-01
                                        </span>
                                    </div>
                                </div>
                                <div className="text-sm flex flex-row flex-1 gap-2 text-gray-600 justify-evenly">
                                    <p>
                                        Follwers: 0
                                    </p>
                                    <p>
                                        Follwing: 0

                                    </p>
                                </div>
                                {isOwnProfile && (
                                    <button onClick={handleEditOrDone}>
                                        <div className="bg-black text-white hover:bg-gray-600 px-4 py-2 rounded ">
                                            <div className="flex items-center gap-2">
                                                <p>{editOnClick ? "Done" : "Edit"}</p>
                                                {editOnClick ? (
                                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-4">
                                                        <path strokeLinecap="round" strokeLinejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                                                    </svg>
                                                ) : (
                                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-4">
                                                        <path strokeLinecap="round" strokeLinejoin="round" d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L6.832 19.82a4.5 4.5 0 0 1-1.897 1.13l-2.685.8.8-2.685a4.5 4.5 0 0 1 1.13-1.897L16.863 4.487Zm0 0L19.5 7.125" />
                                                    </svg>
                                                )}
                                            </div>
                                        </div>
                                    </button>
                                )}
                            </div>
                            <div className="flex items-center gap-4 p-4">
                                <textarea
                                    className="flex-1 p-2 rounded-md focus:outline-none focus:ring-2 focus:ring-[#B6B09F] focus:border-[#B6B09F] resize-none"
                                    disabled={!editOnClick}
                                    rows={1}
                                    value={bio}
                                    onChange={e => setBio(e.target.value)}
                                />

                            </div>

                        </div>
                        <Feed/>

                    </div>

                </main>
            </div>
        </Layout>
    );
}