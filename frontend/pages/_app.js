// pages/_app.js
import "@/styles/globals.css"; //
import { AuthProvider } from '@/contexts/AuthContext'; // 引入 AuthProvider

export default function App({ Component, pageProps }) {
  return (
    <AuthProvider> 
      <Component {...pageProps} /> 
    </AuthProvider>
  );
}