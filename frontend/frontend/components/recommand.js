export default function Recommand() {
    return (
        <div className="bg-white p-5 rounded-lg shadow">
        <h3 className="text-lg font-semibold text-gray-700 mb-3">推薦好友</h3>
        <ul className="space-y-2">
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