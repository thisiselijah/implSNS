// components/avatar.js
import { useAuth } from '@/contexts/AuthContext'; // 1. 引入 useAuth
import { useRouter } from 'next/router';   // 1. 引入 useRouter

export default function Avatar( props ) {
  const authContext = useAuth(); // 2. 獲取 AuthContext 實例
  const router = useRouter();   // 2. 獲取 router 實例

  let username = props.username ? props.username : "Default User";

  const handleLogout = async () => {
    console.log("Logging out...");
    const storedToken = localStorage.getItem('jwtToken'); // 在函數執行時獲取最新的 token

    if (!storedToken) {
      console.error("No token found for logout.");
      // 可能需要提示用戶或直接嘗試清除前端狀態並跳轉
      if (authContext && typeof authContext.logout === 'function') {
        authContext.logout();
      }
      router.push('/');
      return;
    }

    try {
      const response = await fetch('http://localhost:8080/api/v1/auth/logout', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${storedToken}`, // 3. 修正拼寫錯誤
          // 'Content-Type': 'application/json', // 根據後端 API 是否需要決定是否添加
        },
      });

      if (!response.ok) {
        // 嘗試解析錯誤訊息
        let errorData;
        try {
          errorData = await response.json();
        } catch (e) {
          // 如果回應不是 JSON 或解析失敗
          errorData = { message: response.statusText };
        }
        console.error("Logout failed on server:", errorData.message || response.statusText);
        alert(`登出失敗：${errorData.message || response.statusText}`); // 4. 給用戶反饋
        if (authContext && typeof authContext.logout === 'function') {
          authContext.logout(); // 清除前端狀態
        }
        router.push('/'); // 導向到首頁或登入頁
        return;
      }

      // 假設後端成功登出後可能返回一些信息，如果沒有則 response.json() 可能會出錯，取決於後端實現
      try {
        const responseData = await response.json(); // 5. 使用 const/let 聲明
        console.log("Logout successful from server:", responseData);
      } catch (e) {
        console.log("Logout successful (no JSON response or parse error, which might be OK).");
        // 如果後端登出成功時沒有返回 JSON body (例如返回 204 No Content)，這裡會出錯，但登出本身是成功的
      }

      if (authContext && typeof authContext.logout === 'function') {
        authContext.logout(); // 清除 Context 和 localStorage 中的 token
      }
      alert("您已成功登出。"); // 4. 給用戶反饋
      router.push('/'); // 導向到首頁或登入頁

    } catch (error) {
      console.error("An error occurred during logout:", error);
      alert(`登出過程中發生錯誤：${error.message}`); // 4. 給用戶反饋
      if (authContext && typeof authContext.logout === 'function') {
        authContext.logout();
      }
      router.push('/');
    }
  };

  return (
    <>
      <div className="bg-white p-4 rounded-lg shadow text-center space-y-4">
        <div className="flex justify-center">
          <img
            alt={username}
            src="https://images.unsplash.com/photo-1491528323818-fdd1faba62cc?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
            className="inline-block size-16 rounded-full ring-2 ring-offset-2 ring-[#B6B09F]"
          />
        </div>
        <h3 className="text-xl font-semibold text-gray-800">{username}</h3>
        <button
          onClick={handleLogout}
          className="w-full flex items-center justify-center space-x-2 bg-black hover:bg-gray-600 text-white font-medium py-2.5 px-4 rounded-md transition-colors duration-150"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth={1.5}
            stroke="currentColor"
            className="size-5"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M15.75 9V5.25A2.25 2.25 0 0 0 13.5 3h-6a2.25 2.25 0 0 0-2.25 2.25v13.5A2.25 2.25 0 0 0 7.5 21h6a2.25 2.25 0 0 0 2.25-2.25V15M12 9l-3 3m0 0 3 3m-3-3h12.75"
            />
          </svg>
          <span>Log Out</span>
        </button>
      </div>
    </>
  );
}