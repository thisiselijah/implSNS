// components/leftpanel.js
import Link from 'next/link'; // 建議使用 Next.js 的 Link 組件進行內部導航

export default function LeftPanel() {
  const navItems = [
    {
      href: "/",
      icon: (
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6 flex-shrink-0">
          <path strokeLinecap="round" strokeLinejoin="round" d="m2.25 12 8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25" />
        </svg>
      ),
      text: "Home",
    },
    {
      href: "/notifications", // 假設通知頁面的路徑
      icon: (
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6 flex-shrink-0">
          <path strokeLinecap="round" strokeLinejoin="round" d="M14.857 17.082a23.848 23.848 0 0 0 5.454-1.31A8.967 8.967 0 0 1 18 9.75V9A6 6 0 0 0 6 9v.75a8.967 8.967 0 0 1-2.312 6.022c1.733.64 3.56 1.085 5.455 1.31m5.714 0a24.255 24.255 0 0 1-5.714 0m5.714 0a3 3 0 1 1-5.714 0" />
        </svg>
      ),
      text: "Notifications",
    },
    {
      href: "/settings", // 假設設定頁面的路徑
      icon: (
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6 flex-shrink-0">
          <path strokeLinecap="round" strokeLinejoin="round" d="M10.5 6h9.75M10.5 6a1.5 1.5 0 1 1-3 0m3 0a1.5 1.5 0 1 0-3 0M3.75 6H7.5m3 12h9.75m-9.75 0a1.5 1.5 0 0 1-3 0m3 0a1.5 1.5 0 0 0-3 0m-3.75 0H7.5m9-6h3.75m-3.75 0a1.5 1.5 0 0 1-3 0m3 0a1.5 1.5 0 0 0-3 0m-9.75 0h9.75" />
        </svg>
      ),
      text: "Settings",
    },
  ];

  return (
    <aside className="hidden lg:block lg:col-span-3">
      <div className="sticky top-70 space-y-2 p-1">
        <nav className="bg-white p-3 rounded-lg shadow">
          <ul className="space-y-1">
            {navItems.map((item) => (
              <li key={item.text}>
                {/* 移除 legacyBehavior, 將 className 直接應用於 Link */}
                <Link
                  href={item.href}
                  className="flex items-center space-x-3 p-3 rounded-md text-black hover:bg-gray-50 hover:text-[#B6B09F] transition-colors duration-150 group focus:outline-none focus:ring-2 focus:ring-indigo-500"
                >
                  {/* Link 組件現在直接包裹內容 */}
                  {item.icon}
                  <span className="font-medium group-hover:font-semibold">{item.text}</span>
                </Link>
              </li>
            ))}
          </ul>
        </nav>
      </div>
    </aside>
  );
}