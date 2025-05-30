import Layout from "@/components/layout";
import { PostCard } from "@/components/card";
import Avatar from "@/components/avatar";
import Recommand from "@/components/recommand";
// import Image from "next/image";
import {PostsNavbar} from "@/components/navbar"; // 假設你有一個 PostsNavbar 組件

export default function Post() {
  return (
    <Layout pageTitle="View Post"> {/* 假設 Layout 組件可以接受 pageTitle prop */}
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
        <aside className="hidden lg:block lg:col-span-3">
          <div className="sticky top-70 space-y-6"> {/* sticky top-20 讓它在捲動時固定在頂部某個位置 */}
            <div className="bg-white p-5 rounded-lg shadow">
              <ul className="space-y-2">
                <li className="flex items-center space-x-2">
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
                    <path strokeLinecap="round" strokeLinejoin="round" d="m2.25 12 8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25" />
                  </svg>
                  <a href="/" className="text-[#0D0D0D] hover:text-[#0D0D0D] hover:underline">Home</a>
                </li>
                <li className="flex items-center space-x-2">

                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
                    <path strokeLinecap="round" strokeLinejoin="round" d="M14.857 17.082a23.848 23.848 0 0 0 5.454-1.31A8.967 8.967 0 0 1 18 9.75V9A6 6 0 0 0 6 9v.75a8.967 8.967 0 0 1-2.312 6.022c1.733.64 3.56 1.085 5.455 1.31m5.714 0a24.255 24.255 0 0 1-5.714 0m5.714 0a3 3 0 1 1-5.714 0" />
                  </svg>

                  <a href="#" className="text-[#0D0D0D] hover:text-[#0D0D0D] hover:underline">Notifications</a>
                </li>
                <li className="flex items-center space-x-2">
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
                    <path strokeLinecap="round" strokeLinejoin="round" d="M10.5 6h9.75M10.5 6a1.5 1.5 0 1 1-3 0m3 0a1.5 1.5 0 1 0-3 0M3.75 6H7.5m3 12h9.75m-9.75 0a1.5 1.5 0 0 1-3 0m3 0a1.5 1.5 0 0 0-3 0m-3.75 0H7.5m9-6h3.75m-3.75 0a1.5 1.5 0 0 1-3 0m3 0a1.5 1.5 0 0 0-3 0m-9.75 0h9.75" />
                  </svg>

                  <a href="#" className="text-[#0D0D0D] hover:text-[#0D0D0D] hover:underline">Settings</a>
                </li>
              </ul>
            </div>
            {/* 你可以在這裡加入更多左側欄的內容 */}
          </div>
        </aside>

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
          <div className="sticky top-20 space-y-6"> {/* sticky top-20 讓它在捲動時固定在頂部某個位置 */}
            <Avatar />
            <div className="bg-white p-5 rounded-lg shadow">
              <h3 className="text-lg font-semibold text-gray-700 mb-3">搜尋</h3>
              <div className="flex flex-row items-center gap-2"> {/* 1. 使用 gap-2 設定間距，items-center 垂直居中 */}
                <input
                  type="text"
                  placeholder="搜尋貼文、用戶..."
                  className="flex-1 p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#1F98A6] focus:border-[#1F98A6]"
                // 2. flex-1 讓輸入框彈性增長
                //    focus:outline-none 移除預設 outline，改用 ring 和 border
                />
                <button
                  className="p-2 bg-[#000000] hover:bg-[000000] text-white font-semibold rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors duration-150"
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