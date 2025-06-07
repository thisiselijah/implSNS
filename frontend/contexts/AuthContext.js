// contexts/AuthContext.js
import React, { createContext, useContext, useState, useEffect } from "react";
import { useRouter } from "next/router";

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  // 狀態從 token 改為 user，更符合實際情況
  const [user, setUser] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const router = useRouter();

  useEffect(() => {
    // 組件首次掛載時，嘗試透過 API 獲取當前使用者資訊來驗證 session
    const checkUserStatus = async () => {
      try {
        // 假設你有一個 API endpoint 用於獲取當前登入用戶的資訊
        // 瀏覽器會自動攜帶 httpOnly cookie
        const response = await fetch(
          `${process.env.NEXT_PUBLIC_API_BASE_URL}${process.env.NEXT_PUBLIC_STATUS_API}`,
          {
            method: "GET",
            headers: {
              "Content-Type": "application/json",
            },
            credentials: "include", // 非常重要！確保瀏覽器發送 cookie
          }
        );

        if (response.ok) {
          const userData = await response.json();
          setUser(userData); // 登入成功，設置使用者資料
        } else {
          // 例如 401 Unauthorized，表示沒有有效的 cookie 或 session
          setUser(null);
        }
      } catch (error) {
        console.error("Failed to fetch user status:", error);
        setUser(null); // 發生錯誤，視為未登入
      } finally {
        setIsLoading(false);
      }
    };

    checkUserStatus();
  }, []);

  // login 函式現在接收使用者資料，而不是 token
  const login = (userData) => {
    setUser(userData);
    router.push("/posts"); // 登入成功後導向 /posts
  };

  // logout 函式需要呼叫後端 API 來清除 cookie
  const logout = async () => {
    try {
      // 假設你有一個登出的 API endpoint
      await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/v1/auth/logout`,
        {
          method: "POST",
          credentials: "include", // 確保請求帶上 cookie
        }
      );
    } catch (error) {
      console.error("Logout request failed:", error);
    } finally {
      setUser(null); // 無論如何都在前端清除使用者狀態
      // router.push('/'); // 可以選擇重定向
    }
  };

  const isAuthenticated = !!user;

  // 將 user 和新的 login/logout 函式傳遞下去
  return (
    <AuthContext.Provider
      value={{ user, isAuthenticated, login, logout, isLoading }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
