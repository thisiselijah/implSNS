// contexts/AuthContext.js
import React, { createContext, useContext, useState, useEffect } from 'react';
import { useRouter } from 'next/router';

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [token, setToken] = useState(null);
  const [isLoading, setIsLoading] = useState(true); // 用於處理初始 Token 加載狀態
  const router = useRouter();

  useEffect(() => {
    // 組件首次掛載時，嘗試從 localStorage 讀取 Token
    try {
      const storedToken = localStorage.getItem('jwtToken');
      if (storedToken) {
        setToken(storedToken);
      }
    } catch (error) {
      console.error("Could not access localStorage:", error);
      // 在某些環境下 (例如 SSR 期間的某些情況或隱私模式) localStorage 可能不可用
    }
    setIsLoading(false);
  }, []);

  const login = (newToken) => {
    try {
      localStorage.setItem('jwtToken', newToken);
      setToken(newToken);
    } catch (error) {
      console.error("Could not access localStorage:", error);
    }
  };

  const logout = () => {
    try {
      localStorage.removeItem('jwtToken');
    } catch (error) {
      console.error("Could not access localStorage:", error);
    }
    setToken(null);
    // 您可以選擇在這裡將用戶重定向到登入頁或首頁
    // router.push('/'); // 例如，登出後返回首頁
  };

  const isAuthenticated = !!token;

  return (
    <AuthContext.Provider value={{ token, isAuthenticated, login, logout, isLoading }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};