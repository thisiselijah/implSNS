// pages/posts/post.js
import Layout from "@/components/layout"; //
import { PostCard } from "@/components/card"; //
import Avatar from "@/components/avatar"; //
import Recommand from "@/components/recommand"; //
import { PostsNavbar } from "@/components/navbar"; //
import { useAuth } from '@/contexts/AuthContext'; // 引入 useAuth
import { useRouter } from 'next/router';
import { useEffect } from 'react';

export default function Post() {
  const { isAuthenticated, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    // 等待 AuthContext 完成初始 Token 加載
    if (!isLoading) {
      if (!isAuthenticated) {
        // 如果用戶未認證，則重定向
        // 假設您的登入是通過首頁的 Auth 組件（模態框）進行的
        alert('您需要登入才能查看此頁面。將返回首頁。');
        router.replace('/'); // 重定向到首頁
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

  // 用戶已認證，渲染實際的頁面內容
  return (
    <Layout pageTitle="View Post"> {/* */}
      <PostsNavbar /> {/* */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 px-4 py-8"> {/* */}

        <aside className="hidden lg:block lg:col-span-3"> {/* */}
          <div className="sticky top-70 space-y-6"> {/* */}
            <div className="bg-white p-5 rounded-lg shadow"> {/* */}
              <ul className="space-y-2"> {/* */}
                <li className="flex items-center space-x-2"> {/* */}
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6"> {/* */}
                    <path strokeLinecap="round" strokeLinejoin="round" d="m2.25 12 8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25" /> {/* */}
                  </svg>
                  <a href="/" className="text-[#0D0D0D] hover:text-[#0D0D0D] hover:underline">Home</a> {/* */}
                </li>
                {/* ...其他側邊欄連結... */}
              </ul>
            </div>
          </div>
        </aside>

        <div className="flex flex-col lg:col-span-6 md:col-span-8 col-span-12 space-y-2 "> {/* */}
          <PostCard /> {/* */}
          <PostCard /> {/* */}
          <PostCard /> {/* */}
        </div>

        <aside className="hidden md:block md:col-span-4 lg:col-span-3"> {/* */}
          <div className="sticky top-20 space-y-6"> {/* */}
            <Avatar /> {/* */}
            <div className="bg-white p-5 rounded-lg shadow"> {/* */}
              {/* ...右側邊欄內容... */}
            </div>
            <Recommand /> {/* */}
          </div>
        </aside>
      </div>
    </Layout>
  );
}