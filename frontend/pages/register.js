import Layout from "@/components/layout";
import { useRouter } from "next/router";
import { useState } from "react";

export default function Register() {
    const router = useRouter();
    const [loading, setLoading] = useState(false);
    const [errorMsg, setErrorMsg] = useState(""); // 新增錯誤訊息狀態

    const handleFormSubmit = async (event) => {
        event.preventDefault();
        setErrorMsg(""); // 清除前一次錯誤

        const formData = new FormData(event.target);
        const data = Object.fromEntries(formData.entries());

        // 前端驗證
        if (!data.username || data.username.length < 3) {
            setErrorMsg("使用者名稱至少需 3 個字元");
            return;
        }
        if (!data.email || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(data.email)) {
            setErrorMsg("請輸入正確的電子郵件格式");
            return;
        }
        if (!data.password || data.password.length < 8) {
            setErrorMsg("密碼至少需 8 個字元");
            return;
        }

        setLoading(true);
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
                setErrorMsg(errorData.error || "註冊失敗，請稍後再試。");
            } else {
                setErrorMsg("");
                alert("註冊成功！請登入。");
                // 可以在這裡重定向到登入頁面
            }
        } catch (error) {
            setErrorMsg("註冊過程中發生錯誤，請稍後再試。");
        } finally {
            setLoading(false);
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
                            {errorMsg && (
                                <div className="text-red-500 text-center">{errorMsg}</div>
                            )}
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