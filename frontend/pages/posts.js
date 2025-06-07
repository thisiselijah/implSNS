// pages/posts/post.js
import Layout from "@/components/layout";
import Avatar from "@/components/avatar";
import Recommand from "@/components/recommand";
import Feed from "@/components/feed";
import CreatePost from "@/components/createpost";
import { PostsNavbar } from "@/components/navbar";
import { useAuth } from "@/contexts/AuthContext";
import { useRouter } from "next/router";
import { use, useEffect, useState } from "react";
import LeftPanel from "@/components/leftpanel";

export async function getServerSideProps(context) {
  const { req } = context;
  // 從 cookie 取得 jwt_token
  const jwtToken = req.cookies.jwt_token;

  // 先檢查 jwt_token 是否存在
  if (!jwtToken) {
    return {
      redirect: {
        destination: "/",
        permanent: false,
      },
    };
  }

  // 透過 /auth/status 取得 userId
  let userId = null;
  try {
    const statusRes = await fetch(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/v1/auth/status`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          cookie: `jwt_token=${jwtToken}`,
        },
      }
    );
    if (statusRes.ok) {
      const statusData = await statusRes.json();
      userId = statusData.userID;
    }
  } catch (e) {
    // ignore
  }

  if (!userId) {
    return {
      redirect: {
        destination: "/",
        permanent: false,
      },
    };
  }

  try {
    // 取得貼文動態
    const feedResponse = await fetch(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}${process.env.NEXT_PUBLIC_POSTS_FEED_API}${userId}`,
      {
        method: "GET",
        headers: {
          cookie: `jwt_token=${jwtToken}`,
        },
      }
    );

    if (!feedResponse.ok) {
      throw new Error(`無法載入動態 (Code: ${feedResponse.status})`);
    }
    const initialFeedData = await feedResponse.json();

    const profileResponse = await fetch(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}${process.env.NEXT_PUBLIC_PROFILE_API}${userId}`,
      {
        method: "GET",
        headers: {
          cookie: `jwt_token=${jwtToken}`,
        },
      }
    );

    if (!profileResponse.ok) {
      throw new Error(`無法載入使用者資料 (Code: ${profileResponse.status})`);
    }
    const profileData = await profileResponse.json();
    if (!profileData || !profileData.user_id) {
      throw new Error("個人資料格式不正確");
    }

    return {
      props: {
        userId,
        initialFeedData,
        avatar_url: profileData.avatar_url || null,
        username: profileData.username || null,
        error: null,
      },
    };
  } catch (error) {
    console.error("SSR: Error in getServerSideProps:", error);
    return {
      props: {
        userId: null,
        initialFeedData: null,
        avatar_url: null,
        username: null,
        error: error.message || "伺服器發生錯誤，請稍後再試。",
      },
    };
  }
}

export default function Posts({
  userId,
  initialFeedData,
  avatar_url,
  username,
  error,
}) {
  // 客戶端的 useAuth hook 仍然非常重要！
  const { isAuthenticated, isLoading: authIsLoading } = useAuth();
  const router = useRouter();

  const [isCreatePostOpen, setIsCreatePostOpen] = useState(false);
  const handleCreatePostClick = () => setIsCreatePostOpen(true);
  const handleCloseCreatePost = () => setIsCreatePostOpen(false);

  useEffect(() => {
    if (!authIsLoading && !isAuthenticated && !error) {
      alert("您的登入已逾期或已登出，將返回首頁。");
      router.replace("/");
    }
  }, [isAuthenticated, authIsLoading, router, error]);

  // 處理 SSR 期間發生的錯誤
  if (error) {
    return (
      <Layout pageTitle="錯誤">
        <PostsNavbar />
        <div className="flex flex-col gap-2 justify-center items-center min-h-screen">
          <p className="text-red-500 text-lg">{error}</p>
          <button
            onClick={() => router.push("/")}
            className="mt-4 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            返回首頁
          </button>
        </div>
      </Layout>
    );
  }

  // 初始載入狀態 (主要由 getServerSideProps 處理，客戶端載入是 useAuth 的狀態)
  // 如果 initialFeedData 存在，表示 SSR 成功，理論上 isAuthenticated 也應為 true
  // 如果 authIsLoading 為 true，可能 useAuth 正在進行一些初始化檢查
  if (authIsLoading && !initialFeedData) {
    return (
      <Layout pageTitle="載入中...">
        <PostsNavbar />
        <div className="flex flex-col gap-2 justify-center items-center min-h-screen">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
          <p className="text-gray-500 ml-4">載入中，請稍候...</p>
        </div>
      </Layout>
    );
  }

  const feedData = (initialFeedData) => {
    if (!initialFeedData || !Array.isArray(initialFeedData)) {
      console.warn("SSR: Initial feed data is not an array or is null.");
      return [];
    }
    return initialFeedData.map((item) => ({
      post_id: item.post_id,
      author_id: item.author_id,
      author_name: item.author_name,
      content: item.content,
      media: Array.isArray(item.media) ? item.media : [],
      tags: Array.isArray(item.tags) ? item.tags : [],
      location: item.location || null,
      like_count: typeof item.like_count === "number" ? item.like_count : 0,
      comment_count:
        typeof item.comment_count === "number" ? item.comment_count : 0,
      created_at: item.created_at ? new Date(item.created_at) : null,
      updated_at: item.updated_at ? new Date(item.updated_at) : null,
    }));
  };

  // 主要內容渲染
  return (
    <Layout pageTitle="View Post">
      <PostsNavbar userId={userId} />
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-4 px-4 min-h-screen">
        <LeftPanel />{" "}
        {/* LeftPanel 可能也需要 SSR 資料，如果其內容是使用者特定的 */}
        <main className="flex flex-col lg:col-span-6 md:col-span-8 col-span-12">
          <div className="bg-white p-0.5 rounded-lg ">
            <div className="p-5 flex flex-row items-center">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
                stroke="currentColor"
                className="size-6"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0 1 15.75 21H5.25A2.25 2.25 0 0 1 3 18.75V8.25A2.25 2.25 0 0 1 5.25 6H10"
                />
              </svg>
              <button
                onClick={handleCreatePostClick}
                className="text-left italic text-gray-600 flex-1 ml-2 p-2 rounded-md focus:outline-none focus:ring-2 focus:ring-[#B6B09F] focus:border-[#B6B09F]"
              >
                Got something to share today?
              </button>
            </div>
            <Feed feedData={feedData(initialFeedData)} />
          </div>
        </main>
        <aside className="hidden md:block md:col-span-4 lg:col-span-3">
          <div className="flex flex-col gap-2 sticky top-30">
            {/* Avatar 可能需要從 authContext 或 serverRenderedUserId 獲取用戶資訊 */}
            <Avatar
              router={router}
              avatar_url={avatar_url}
              username={username}
            />
            <div className="flex flex-col gap-2 bg-white rounded-lg shadow">
              <div className="flex flex-row items-center gap-2 p-4 border-b border-gray-200">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  strokeWidth={1.5}
                  stroke="currentColor"
                  className="size-6"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z"
                  />
                </svg>
                <h3 className="text-lg font-semibold text-gray-700">搜尋</h3>
              </div>

              <div className="flex flex-row items-center p-4 gap-2">
                <input
                  type="text"
                  placeholder="搜尋貼文、用戶..."
                  className="flex-1 p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#B6B09F] focus:border-[#B6B09F]"
                />
                <button className="p-2 bg-[#000000] hover:bg-gray-600 text-white font-semibold rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors duration-150">
                  搜尋
                </button>
              </div>
            </div>
            <Recommand /> {/* Recommand 可能也需要 SSR 資料 */}
          </div>
        </aside>
      </div>

      {isCreatePostOpen && (
        <div className="fixed inset-0 bg-black/85 flex items-center justify-center z-50 p-4">
          <div className="bg-transparent rounded-lg relative w-full max-w-xl">
            <CreatePost
              onClose={handleCloseCreatePost}
              avatar_url={avatar_url}
            />
          </div>
        </div>
      )}
    </Layout>
  );
}
