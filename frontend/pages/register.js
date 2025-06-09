import Layout from "@/components/layout";
import { useRouter } from "next/router";
import { useState } from "react";

export default function Register() {
    const router = useRouter();
    const [loading, setLoading] = useState(false);
    const handleFormSubmit = async (event) => {
        event.preventDefault();

        const formData = new FormData(event.target);
        const data = Object.fromEntries(formData.entries());
        console.log("註冊資料：", data);
        

        try {
            const response = await fetch(
                process.env.NEXT_PUBLIC_API_BASE_URL + process.env.NEXT_PUBLIC_REGISTER_API, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(data),
            });

            if (!response.ok) {
                const errorData = await response.json();
                alert(errorData.error || "註冊失敗，請稍後再試。");
            } else {
                const result = await response.json();
                alert("註冊成功！請登入。");
                // 可以在這裡重定向到登入頁面或其他操作
            }
        } catch (error) {
            console.error("註冊過程中發生錯誤：", error);
            alert("註冊過程中發生錯誤，請稍後再試。");
        }
    };

    return (
                <Layout>
          <div className="flex items-center justify-center min-h-screen">
            <main className="flex flex-col w-full max-w-md">
              <div className="bg-white p-0.5 rounded-lg shadow">
                <div className="flex flex-col p-4 border-b border-gray-200">
                  <h1 className="text-2xl font-bold text-center">註冊</h1>
                </div>
                <form className="flex flex-col gap-4 p-6" onSubmit={handleFormSubmit}>
                  <input
                    name="username"
                    type="text"
                    placeholder="使用者名稱"
                    required
                    className="p-2 border rounded"
                  />
                  <input
                    name="email"
                    type="email"
                    placeholder="電子郵件"
                    required
                    className="p-2 border rounded"
                  />
                  <input
                    name="password"
                    type="password"
                    placeholder="密碼"
                    required
                    className="p-2 border rounded"
                  />
                  <button
                    type="submit"
                    className="bg-black text-white py-2 rounded hover:bg-gray-700"
                    disabled={loading}
                  >
                    {loading ? "註冊中..." : "註冊"}
                  </button>
                </form>
              </div>
            </main>
          </div>
        </Layout>
    );
}