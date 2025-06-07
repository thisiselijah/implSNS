// auth.js
import { useAuth } from '@/contexts/AuthContext'; // 引入 useAuth hook

export default function Auth({ onClose }) { // 不再需要 onLoginSuccess
  const { login } = useAuth(); // 從 context 取得 login 方法

  const handleFormSubmit = async (event) => {
    event.preventDefault();

    console.log("Form submitted in Auth component");
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}${process.env.NEXT_PUBLIC_LOGIN_API}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email: event.target.email.value,
          password: event.target.password.value,
        }),
        credentials: 'include', // 重要！確保瀏覽器接收 Set-Cookie header
      });

      if (!response.ok) {
        // 錯誤處理邏輯保持不變
        const errData = await response.json();
        throw new Error(errData.error || "Login failed.");
      }
      
      // 登入成功，後端已設置 cookie
      const userData = await response.json();
      console.log("Login successful, user data:", userData);

      login(userData);

      if (onClose) {
        onClose();
      }

    } catch (error) {
      console.error("There was a problem with the login request:", error);
      alert(error.message || "Login failed. Please check your credentials and try again.");
    }
  };

  return (
    // ... JSX 保持不變 ...
    // [JSX from file 1]
    <div className="bg-white p-6 sm:p-8 rounded-lg shadow-xl w-full max-w-md transform transition-all">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold tracking-tight text-gray-900">
          Sign in to your account
        </h2>
        {onClose && (
          <button
            onClick={onClose}
            type="button"
            className="text-gray-400 hover:text-gray-600 p-1 rounded-full hover:bg-gray-100"
            aria-label="Close"
          >
            <svg
              className="h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        )}
      </div>
      <div className="sm:mx-auto sm:w-full sm:max-w-sm mb-6">
        <img
          alt="Your Company"
          src="https://tailwindcss.com/plus-assets/img/logos/mark.svg?color=indigo&shade=600"
          className="mx-auto h-10 w-auto"
        />
      </div>
      <form onSubmit={handleFormSubmit} className="space-y-6">
        <div>
          <label
            htmlFor="email"
            className="block text-sm font-medium leading-6 text-gray-900"
          >
            Email address
          </label>
          <div className="mt-2">
            <input
              id="email"
              name="email"
              type="email"
              required
              autoComplete="email"
              className="block w-full rounded-md border-0 py-1.5 px-3 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
            />
          </div>
        </div>
        <div>
          <div className="flex items-center justify-between">
            <label
              htmlFor="password"
              className="block text-sm font-medium leading-6 text-gray-900"
            >
              Password
            </label>
            <div className="text-sm">
              <a
                href="#"
                className="font-semibold text-indigo-600 hover:text-indigo-500"
              >
                Forgot password?
              </a>
            </div>
          </div>
          <div className="mt-2">
            <input
              id="password"
              name="password"
              type="password"
              required
              autoComplete="current-password"
              className="block w-full rounded-md border-0 py-1.5 px-3 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
            />
          </div>
        </div>
        <div>
          <button
            type="submit"
            className="flex w-full justify-center rounded-md bg-indigo-600 px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
          >
            Sign in
          </button>
        </div>
      </form>
      <p className="mt-10 text-center text-sm text-gray-500">
        Not a member?
        <a
          href="/register"
          className="font-semibold leading-6 text-indigo-600 hover:text-indigo-500"
        >
          <span aria-hidden="true"> </span>
          Create a new account
        </a>
      </p>
    </div>
  );
}