// auth.js
export default function Auth({ onClose, onLoginSuccess }) {

  const handleFormSubmit = async (event) => {
    event.preventDefault();

    console.log("Form submitted in Auth component");
    try {
      const response = await fetch("http://localhost:8080/api/v1/auth/login", {
        //
        method: "POST", //
        headers: {
          //
          "Content-Type": "application/json", //
        },
        body: JSON.stringify({
          //
          email: event.target.email.value, //
          password: event.target.password.value, //
        }),
      });

      if (!response.ok) {
        //
        // 如果 HTTP 狀態碼不是 2xx，則認為是錯誤
        try {
          const errData = await response.json();
          // 嘗試解析錯誤回應的 JSON 主體
          throw new Error(errData.message || "Network response was not ok");
        } catch {
          // 如果錯誤回應不是 JSON 或解析失敗
          throw new Error(
            "Network response was not ok and error body is not valid JSON"
          );
        }
      }

      const data = await response.json(); //
      console.log("Login successful:", data); //

      const jwtToken = data.token || data.accessToken;

      if (jwtToken) {
        localStorage.setItem("jwtToken", jwtToken);
        console.log("Token saved to localStorage:", jwtToken);

        if (onLoginSuccess) {
          onLoginSuccess(jwtToken); // 傳遞 token 給父組件
        }
        if (onClose) {
          // 關閉 Modal 的邏輯可以保留或也由父組件處理
          onClose();
        }
      } else {
        console.warn("JWT Token not found in the login response data.");
        // 可以在這裡處理 Token 未返回的情況，例如顯示一個通用錯誤訊息
        alert(
          "Login successful, but no token was received. Please contact support."
        );
      }
    } catch (error) {
      //
      console.error("There was a problem with the login request:", error); //
      // 在這裡可以處理登入失敗的邏輯，例如顯示錯誤訊息給用戶
      // 例如: alert(error.message || 'Login failed. Please check your credentials.');
      // 實際應用中，你可能會將錯誤訊息顯示在表單附近而不是用 alert
      const errorMessage =
        error.message ||
        "Login failed. Please check your credentials and try again.";
      // 假設你有一個顯示錯誤訊息的 UI 元素
      // displayLoginError(errorMessage);
      alert(errorMessage); // 暫時用 alert 顯示
    }
  };

  return (
    // ... 其餘 JSX 代碼保持不變 ...
    // [JSX from file 1]
    // 3. 這個 div 就是「小視窗」本身，設定它的外觀
    <div className="bg-white p-6 sm:p-8 rounded-lg shadow-xl w-full max-w-md transform transition-all">
      {" "}
      {/* */}
      {/* bg-white: 白色背景
        p-6 sm:p-8: 內邊距，小螢幕以上有更大內邊距
        rounded-lg: 較大的圓角
        shadow-xl: 較深的陰影，使其有立體感
        w-full max-w-md: 寬度100%，但最大寬度限制在中等尺寸 (md)，避免在大螢幕上過寬
        transform transition-all: (可選) 為將來的動畫效果做準備
      */}
      {/* 4. 模態框的標題和關閉按鈕 */}
      <div className="flex justify-between items-center mb-6">
        {" "}
        {/* */}
        <h2 className="text-2xl font-bold tracking-tight text-gray-900">
          {" "}
          {/* */}
          Sign in to your account
        </h2>
        {onClose && ( // 只有當 onClose prop 被傳遞時才顯示關閉按鈕
          <button
            onClick={onClose} //
            type="button" // 避免觸發表單提交
            className="text-gray-400 hover:text-gray-600 p-1 rounded-full hover:bg-gray-100" //
            aria-label="Close" //
          >
            <svg
              className="h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              {" "}
              {/* */}
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M6 18L18 6M6 6l12 12"
              />{" "}
              {/* */}
            </svg>
          </button>
        )}
      </div>
      <div className="sm:mx-auto sm:w-full sm:max-w-sm mb-6">
        {" "}
        {/* */}
        <img
          alt="Your Company" //
          src="https://tailwindcss.com/plus-assets/img/logos/mark.svg?color=indigo&shade=600" //
          className="mx-auto h-10 w-auto" //
        />
      </div>
      {/* 表單本身 - 不需要外層的 sm:mx-auto sm:w-full sm:max-w-sm，因為父 div 已經限制了寬度 */}
      {/* action="#" method="POST" 改為 onSubmit={handleFormSubmit} */}
      <form onSubmit={handleFormSubmit} className="space-y-6">
        {" "}
        {/* */}
        <div>
          <label
            htmlFor="email"
            className="block text-sm font-medium leading-6 text-gray-900"
          >
            {" "}
            {/* */}
            Email address
          </label>
          <div className="mt-2">
            {" "}
            {/* */}
            <input
              id="email" //
              name="email" //
              type="email" //
              required //
              autoComplete="email" //
              className="block w-full rounded-md border-0 py-1.5 px-3 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6" //
            />
          </div>
        </div>
        <div>
          <div className="flex items-center justify-between">
            {" "}
            {/* */}
            <label
              htmlFor="password"
              className="block text-sm font-medium leading-6 text-gray-900"
            >
              {" "}
              {/* */}
              Password
            </label>
            <div className="text-sm">
              {" "}
              {/* */}
              <a
                href="#"
                className="font-semibold text-indigo-600 hover:text-indigo-500"
              >
                {" "}
                {/* */}
                Forgot password?
              </a>
            </div>
          </div>
          <div className="mt-2">
            {" "}
            {/* */}
            <input
              id="password" //
              name="password" //
              type="password" //
              required //
              autoComplete="current-password" //
              className="block w-full rounded-md border-0 py-1.5 px-3 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6" //
            />
          </div>
        </div>
        <div>
          <button
            type="submit" //
            className="flex w-full justify-center rounded-md bg-indigo-600 px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600" //
          >
            Sign in
          </button>
        </div>
      </form>
      <p className="mt-10 text-center text-sm text-gray-500">
        {" "}
        {/* */}
        Not a member? {/* */}
        <a
          href="#"
          className="font-semibold leading-6 text-indigo-600 hover:text-indigo-500"
        >
          {" "}
          {/* */}
          <span aria-hidden="true"> </span> {/* */}
          Create a new account
        </a>
      </p>
    </div>
  );
}
