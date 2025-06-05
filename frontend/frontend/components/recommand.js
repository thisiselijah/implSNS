export default function Recommand() {
    return (
        <div className="bg-white rounded-lg shadow">
            <div className="flex flex-row gap-2 p-4 border-b border-gray-200">
                
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
                    <path strokeLinecap="round" strokeLinejoin="round" d="M18 7.5v3m0 0v3m0-3h3m-3 0h-3m-2.25-4.125a3.375 3.375 0 1 1-6.75 0 3.375 3.375 0 0 1 6.75 0ZM3 19.235v-.11a6.375 6.375 0 0 1 12.75 0v.109A12.318 12.318 0 0 1 9.374 21c-2.331 0-4.512-.645-6.374-1.766Z" />
                </svg>
                <h3 className="text-lg font-semibold text-gray-700">推薦好友</h3>

            </div>

            <ul className="flex flex-col gap-2 p-4">
                <li className="flex items-center space-x-2">
                    <img
                        src="https://images.unsplash.com/photo-1491528323818-fdd1faba62cc?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
                        alt="用戶頭像"
                        className="w-10 h-10 rounded-full"
                    />
                    <span className="text-gray-800">用戶A</span>
                </li>
                <li className="flex items-center space-x-2">
                    <img
                        src="https://images.unsplash.com/photo-1491528323818-fdd1faba62cc?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
                        alt="用戶頭像"
                        className="w-10 h-10 rounded-full"
                    />
                    <span className="text-gray-800">用戶B</span>
                </li>
                <li className="flex items-center space-x-2">
                    <img
                        src="https://images.unsplash.com/photo-1491528323818-fdd1faba62cc?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
                        alt="用戶頭像"
                        className="w-10 h-10 rounded-full"
                    />
                    <span className="text-gray-800">用戶C</span>
                </li>
            </ul>
        </div>
    );
}