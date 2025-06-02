// pages/posts/post.js
import Layout from "@/components/layout"; //
import { PostCard } from "@/components/card"; //
import Avatar from "@/components/avatar"; //
import Recommand from "@/components/recommand"; //
import { PostsNavbar } from "@/components/navbar"; //
import { useAuth } from "@/contexts/AuthContext"; // 引入 useAuth
import { useRouter } from "next/router";
import { useEffect } from "react";
import LeftPanel from "@/components/leftpanel";

export default function Post() {
  const { isAuthenticated, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    // 等待 AuthContext 完成初始 Token 加載
    if (!isLoading) {
      if (!isAuthenticated) {
        // 如果用戶未認證，則重定向
        // 假設您的登入是通過首頁的 Auth 組件（模態框）進行的
        alert("您需要登入才能查看此頁面。將返回首頁。");
        router.replace("/"); // 重定向到首頁
      }
    }
  }, [isAuthenticated, isLoading, router]);

    // 如果正在加載認證狀態或用戶尚未認證（在重定向前），可以顯示加載指示器或 null
  if (isLoading || !isAuthenticated) {
    return (
      <Layout pageTitle="載入中...">
      <PostsNavbar />
      <div className="flex flex-col gap-2 justify-center items-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
        <p className="text-gray-500 ml-4">請先登入...</p>
      </div>
      </Layout>
    );
  }

  return (
    <Layout pageTitle="View Post">
      {" "}
      {/* 假設 Layout 組件可以接受 pageTitle prop */}
      {/* 主要網格容器 (Grid Container)
        - 預設 (最小螢幕): 單欄佈局 (內容會堆疊)
        - md (中等螢幕) 及以上: 變成一個有 12 個可用欄位的網格系統 (grid-cols-12)
        - lg (大型螢幕) 及以上: 我們可以更精確地定義各欄的寬度，例如左右固定，中間彈性
        - gap-6: 設定網格項目之間的間距
        - px-4 py-8: 在網格容器內部加上水平和垂直的內邊距 (padding)
      */}
      <PostsNavbar /> {/* 假設 PostsNavbar 是一個導航欄組件 */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 px-4 py-8">
        {/* --- 左側欄 (其他頁面/導覽) --- */}
        {/* - 在 lg 螢幕及以上，佔據 3/12 的欄位寬度 (lg:col-span-3)
          - 在 lg 以下的螢幕，這個側邊欄會被隱藏 (hidden lg:block) 
            (你也可以選擇讓它在 md 螢幕就顯示，例如用 md:block md:col-span-3)
        */}
        <LeftPanel />
        

        <div className="flex flex-col lg:col-span-6 md:col-span-8 col-span-12 space-y-2 ">
          <PostCard />
          <PostCard />
          <PostCard />
        </div>

        {/* --- 右側欄 (搜尋列、好友建議) --- */}
        {/*
          - 在 lg 螢幕及以上，佔據 3/12 的欄位寬度 (lg:col-span-3)
          - 在 md 螢幕但小於 lg 螢幕時，佔據 4/12 (md:col-span-4)
          - 在 md 以下的螢幕，這個側邊欄會被隱藏 (hidden md:block)
        */}
        <aside className="hidden md:block md:col-span-4 lg:col-span-3">
          <div className="sticky top-40 space-y-6">
            {" "}
            {/* sticky top-20 讓它在捲動時固定在頂部某個位置 */}
            
            <Avatar />

            <div className="bg-white p-5 rounded-lg shadow">
              <h3 className="text-lg font-semibold text-gray-700 mb-3">搜尋</h3>
              <div className="flex flex-row items-center gap-2">
                {" "}
                {/* 1. 使用 gap-2 設定間距，items-center 垂直居中 */}
                <input
                  type="text"
                  placeholder="搜尋貼文、用戶..."
                  className="flex-1 p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#B6B09F] focus:border-[#B6B09F]"
                  // 2. flex-1 讓輸入框彈性增長
                  //    focus:outline-none 移除預設 outline，改用 ring 和 border
                />
                <button
                  className="p-2 bg-[#000000] hover:bg-gray-600 text-white font-semibold rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors duration-150"
                  // 3. 使用 rounded-md 保持一致
                  // 4. 添加更明確的 hover:bg-blue-700 (更深一點的藍色) 和 font-semibold
                  // 5. 添加焦點樣式 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500
                  // 6. 添加 transition-colors 以獲得平滑的背景色過渡效果
                >
                  搜尋
                </button>
              </div>
            </div>
            <Recommand />
          </div>
        </aside>
      </div>
    </Layout>
  );
}
