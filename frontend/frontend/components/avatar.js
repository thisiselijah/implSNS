import { useEffect, useState } from "react";
import { logout } from '@/lib/logout';
import Link from "next/link";

export default function Avatar(props) {
  const authContext = props.authContext;
  const router = props.router;
  const [username, setUsername] = useState("Guest");
  const [user_id, serUser_id] = useState(null);

  useEffect(() => {
    if (typeof window !== "undefined") {
      const name = localStorage.getItem("username") || "Guest";
      setUsername(name);
      const user_id = localStorage.getItem("user_id");
      serUser_id(user_id);
    }
  }, []);

  return (
    <div className="bg-white p-4 rounded-lg shadow text-center space-y-4">
      <div className="flex flex-col gap-2 items-center justify-center">
        <Link href={"/profile/"+user_id} className="mb-2 flex flex-col items-center gap-2">
         <img
          alt={username}
          src="https://images.unsplash.com/photo-1491528323818-fdd1faba62cc?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
          className="inline-block size-16 rounded-full ring-2 ring-offset-2 ring-[#B6B09F]"
        />
        <h3 className="text-xl font-semibold text-gray-800">{username}</h3>
        </Link>
       
      </div>
      
      <button
        onClick={() => logout(authContext, router)}
        className="w-full flex items-center justify-center space-x-2 bg-black hover:bg-gray-600 text-white font-medium py-2.5 px-4 rounded-md transition-colors duration-150"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          strokeWidth={1.5}
          stroke="currentColor"
          className="size-5"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M15.75 9V5.25A2.25 2.25 0 0 0 13.5 3h-6a2.25 2.25 0 0 0-2.25 2.25v13.5A2.25 2.25 0 0 0 7.5 21h6a2.25 2.25 0 0 0 2.25-2.25V15M12 9l-3 3m0 0 3 3m-3-3h12.75"
          />
        </svg>
        <span>Log Out</span>
      </button>
    </div>
  );
}