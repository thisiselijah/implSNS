import { useState } from "react"; // 1. 引入 useState
import Auth from "./auth"; // 假設 Auth.js 在同一目錄
import Link from "next/link";
import { useRouter } from "next/router";
import { useAuth } from '@/contexts/AuthContext';

export function IndexNavbar() {
  // 2. 新增 state 來控制 Auth 組件的顯示狀態
  const [isAuthModalOpen, setIsAuthModalOpen] = useState(false);

  // 3. 處理打開 Auth 模態框的函數
  const handleOpenAuthModal = () => {
    setIsAuthModalOpen(true);
  };

  // 4. 處理關閉 Auth 模態框的函數 (這個函數會傳遞給 Auth 組件)
  const handleCloseAuthModal = () => {
    setIsAuthModalOpen(false);
  };
  const router = useRouter(); // 取得 Next.js 的路由器
  const authContext = useAuth(); // 使用 useAuth 來獲取認證狀態

  const handleLoginSuccess = (token) => {
    // 這裡可以處理登入成功後的邏輯，例如儲存 token 或更新狀態
    authContext.login(token); // 使用 AuthContext 的 login 方法來儲存 token
    router.push("/posts"); // 登入成功後重定向到特定頁面
  }

  return (
    <>
      {/* 使用 Fragment 或 <> </> 來包裹多個頂層元素 */}
      <nav className="bg-black p-3">
        <div className="object-contain flex justify-between items-center">
          <a href="/" className="text-white text-lg font-bold">
            Social Media Project
          </a>


          <div>
            {/* 5. 修改 onClick 事件處理，調用 handleOpenAuthModal */}
            <button
              className="text-black font-bold bg-[#F2F2F2] hover:bg-[F2F2F2] px-4 py-2 rounded"
              onClick={handleOpenAuthModal} // 注意：這裡傳遞的是函數引用，不是函數調用
            >
              Log in
            </button>
          </div>
        </div>
      </nav>
      {/* 6. 條件渲染 Auth 組件作為模態框 */}
      {isAuthModalOpen && (
        <div className="fixed inset-0 bg-black/85 flex items-center justify-center z-50 p-4">
          {/* fixed inset-0: 固定定位，填滿整個視窗
            bg-black bg-opacity-75: 半透明黑色背景遮罩
            flex items-center justify-center: 使用 flex 將內容垂直和水平居中
            z-50: 確保模態框在其他內容之上
            p-4: 為手機等小螢幕設備在邊緣留出一些空間
          */}
          <Auth onClose={handleCloseAuthModal} onLoginSuccess={handleLoginSuccess} />

        </div>
      )}
    </>
  );
}

export function PostsNavbar( {userId} ) {
  return (
    <div className="flex flex-row gap-4 text-center items-center justify-center p-3 py-4.5">
      <Link href={"/posts"}>
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
          <path strokeLinecap="round" strokeLinejoin="round" d="m11.25 9-3 3m0 0 3 3m-3-3h7.5M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
        </svg>
      </Link>
      <p>
        <span className="text-[24px] text-black italic font-bold">Posts</span>
      </p>

      <Link href={"/profile/"+userId}>
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
          <path strokeLinecap="round" strokeLinejoin="round" d="m12.75 15 3-3m0 0-3-3m3 3h-7.5M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
        </svg>
      </Link>

    </div>
  );
}

export function ConstructingNavbar() {
  return (
    <div className="bg-black object-contain p-3 py-4.5">
      <Link href="/" className="text-[18px] text-white hover:underline font-bold">
        Social Media Project
      </Link>
    </div>
  );
}

export function ProfileNavbar() {
  return (
    <div className="flex flex-row gap-4 text-center items-center justify-center p-3 py-4.5">
      <Link href={"/posts"}>
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
          <path strokeLinecap="round" strokeLinejoin="round" d="m11.25 9-3 3m0 0 3 3m-3-3h7.5M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
        </svg>
      </Link>
      <p>
        <span className="text-[24px] text-black italic font-bold">Profile</span>
      </p>
      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
        <path strokeLinecap="round" strokeLinejoin="round" d="m12.75 15 3-3m0 0-3-3m3 3h-7.5M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
      </svg>

    </div>
  );
}
