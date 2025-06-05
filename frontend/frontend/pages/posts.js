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

// 假設這是你後端 API 的基礎 URL
const API_BASE_URL = "http://localhost:8080/api/v1";

export async function getServerSideProps(context) {
  const { req, res } = context;
  let userId = null; // 初始化 userId

  // 1. 從 HTTP-only cookie 中獲取 userID 或 Session ID
  //    你需要知道你的 cookie 名稱，這裡假設是 'auth_token_cookie' 存了 userID，
  //    或者是一個可以換取 userID 的 session_id_cookie。
  //    為了簡化，這裡假設 'user_id_cookie' 直接存了 userID。
  //    在實際應用中，如果存的是 Session ID，你需要一步額外的查詢。
  const userIdFromCookie = req.cookies.user_id; // 請替換成你實際的 cookie 名稱

  if (!userIdFromCookie) {
    // 如果沒有 cookie (表示未登入或 session 失效)，重定向到登入頁面
    console.warn("SSR: No user ID found in cookies, redirecting to login.");
    return {
      redirect: {
        destination: "/", // 你的登入頁面路徑
        permanent: false, // 不是永久重定向
      },
    };
  }
  userId = userIdFromCookie;

  try {
    // 2. 使用 userID 獲取 Feed 數據
    const feedApiUrl = `${API_BASE_URL}/pages/posts/feed/${userId}`;
    const feedResponse = await fetch(feedApiUrl, {
      // 如果 API 需要特定的 headers (例如內部認證 token)，在這裡添加
      // headers: {
      //   'Authorization': `Bearer ${server_side_api_token}` // 舉例
      // }
    });

    if (!feedResponse.ok) {
      // 處理 API 錯誤
      console.error(`SSR: Failed to fetch feed data for user ${userId}. Status: ${feedResponse.status}`);
      if (feedResponse.status === 401 || feedResponse.status === 403) {
        // 權限不足或未授權，可能 cookie 失效，重定向到登入
        return {
          redirect: {
            destination: "/",
            permanent: false,
          },
        };
      }
      // 其他錯誤，可以傳遞一個錯誤狀態給頁面組件
      return { props: { initialFeedData: null, error: `無法載入動態 (Code: ${feedResponse.status})` } };
    }

    const initialFeedData = await feedResponse.json();


    // 3. 將數據作為 props 傳遞給頁面組件
    return {
      props: {
        initialFeedData, // 將 Feed 數據傳遞給組件
        // 你也可以選擇將 userID 傳下去，如果客戶端組件後續操作仍需參考
        // serverRenderedUserId: userId,
      },
    };
  } catch (error) {
    console.error("SSR: Error in getServerSideProps:", error);
    // 發生未知錯誤，可以傳遞一個錯誤狀態
    return {
      props: { initialFeedData: null, error: "伺服器發生錯誤，請稍後再試。" },
    };
  }
}

export default function Posts({ initialFeedData, error /*, serverRenderedUserId */ }) {
  // 客戶端的 useAuth hook 仍然重要，用於管理 UI 狀態、登出、以及 CSR 下的行為
  const { isAuthenticated, isLoading: authIsLoading, user } = useAuth();
  const router = useRouter();
  const authContext = useAuth();

  const [isCreatePostOpen, setIsCreatePostOpen] = useState(false);

  const handleCreatePostClick = () => setIsCreatePostOpen(true);
  const handleCloseCreatePost = () => setIsCreatePostOpen(false);

  // 客戶端身份驗證和重定向邏輯
  // SSR 後，isAuthenticated 應該為 true (因為 getServerSideProps 已經處理了未登入情況)
  // 但這個 useEffect 仍然可以作為一個補充，處理例如 token 在客戶端過期等情況
  useEffect(() => {
    // 如果 SSR 階段出錯導致重定向，這個 effect 可能不會立即執行
    // 主要關注的是 SSR 成功後，客戶端的狀態變化
    if (!authIsLoading && !isAuthenticated && !error) { // 如果沒有 SSR 錯誤，但客戶端認為未驗證
      alert("您的登入已逾期或已登出，將返回首頁。");
      router.replace("/");
    }
  }, [isAuthenticated, authIsLoading, router, error]);

  const [userId, setUserId] = useState(null);

  useEffect(() => {
    setUserId(localStorage.getItem("user_id"));
  }, []);
  



  // 處理 SSR 期間發生的錯誤
  if (error) {
    return (
      <Layout pageTitle="錯誤">
        <PostsNavbar />
        <div className="flex flex-col gap-2 justify-center items-center min-h-screen">
          <p className="text-red-500 text-lg">{error}</p>
          <button onClick={() => router.push('/')} className="mt-4 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600">
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
    return initialFeedData.map(item => ({
      post_id: item.post_id,
      author_name: item.author_name,
      content: item.content,
      media: Array.isArray(item.media) ? item.media : [],
      tags: Array.isArray(item.tags) ? item.tags : [],
      location: item.location || null,
      like_count: typeof item.like_count === 'number' ? item.like_count : 0,
      comment_count: typeof item.comment_count === 'number' ? item.comment_count : 0,
      created_at: item.created_at ? new Date(item.created_at) : null,
      updated_at: item.updated_at ? new Date(item.updated_at) : null,
    }));

  }





  // 主要內容渲染
  return (
    <Layout pageTitle="View Post">
      <PostsNavbar userId={userId}/>
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-4 px-4 min-h-screen">
        <LeftPanel /> {/* LeftPanel 可能也需要 SSR 資料，如果其內容是使用者特定的 */}

        <main className="flex flex-col lg:col-span-6 md:col-span-8 col-span-12">
          <div className="bg-white p-0.5 rounded-lg ">
            <div className="p-5 flex flex-row items-center">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
                <path strokeLinecap="round" strokeLinejoin="round" d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0 1 15.75 21H5.25A2.25 2.25 0 0 1 3 18.75V8.25A2.25 2.25 0 0 1 5.25 6H10" />
              </svg>
              <button
                onClick={handleCreatePostClick}
                className="text-left italic text-gray-600 flex-1 ml-2 p-2 rounded-md focus:outline-none focus:ring-2 focus:ring-[#B6B09F] focus:border-[#B6B09F]"
              >
                Got something to share today?
              </button>
            </div>
            {/* 將從 SSR 獲取的數據傳遞給 Feed 組件 */}
            {/* 確保 Feed 組件能處理 initialData 為 null 或 undefined 的情況 (如果 SSR 獲取失敗但未重定向) */}
            <Feed FeedData={feedData(initialFeedData)} />
          </div>
        </main>

        <aside className="hidden md:block md:col-span-4 lg:col-span-3">
          <div className="flex flex-col gap-2 sticky top-30">
            {/* Avatar 可能需要從 authContext 或 serverRenderedUserId 獲取用戶資訊 */}
            <Avatar router={router} authContext={authContext} />

            <div className="flex flex-col gap-2 bg-white rounded-lg shadow">
              <div className="flex flex-row items-center gap-2 p-4 border-b border-gray-200">
                
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
                  <path strokeLinecap="round" strokeLinejoin="round" d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z" />
                </svg>
                <h3 className="text-lg font-semibold text-gray-700">搜尋</h3>
              </div>

              <div className="flex flex-row items-center p-4 gap-2">
                <input
                  type="text"
                  placeholder="搜尋貼文、用戶..."
                  className="flex-1 p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#B6B09F] focus:border-[#B6B09F]"
                />
                <button
                  className="p-2 bg-[#000000] hover:bg-gray-600 text-white font-semibold rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors duration-150"
                >
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
            <CreatePost onClose={handleCloseCreatePost} />
          </div>
        </div>
      )}
    </Layout>
  );
}