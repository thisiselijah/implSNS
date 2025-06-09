import Image from "next/image";
import Link from "next/link";
import { useState } from "react";

// NavigationCard 組件保持不變
export function NavigationCard({
  title,
  description,
  link,
  bgColor = "bg-indigo-500",
}) {
  return (
    <Link
      href={link}
      className={`block p-6 rounded-lg shadow-lg text-white hover:opacity-90 transition-opacity ${bgColor}`}
    >
      <h3 className="text-xl font-bold mb-2">{title}</h3>
      <p className="text-sm">{description}</p>
    </Link>
  );
}

export function PostCard({ post, authorProfile }) {
  const {
    post_id = post?.post_id || key,
    author_name = authorProfile?.username || "Default User",
    avatar_url = authorProfile?.avatar_url || "/user.png",
    content = "這是貼文的主要內容。Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
    like_count = 14,
    comment_count = 0,
    media = [],
    created_at = "2023-10-01 12:00:00",
  } = post || {};

  const [likes, setLikes] = useState(like_count);
  // isLiked
  const [isLiked, setIsLiked] = useState(post.isLiked); // 可根據需求實作

  // 處理 media 圖片
  const imageSrcs = Array.isArray(media)
    ? media.filter((m) => m.Type === "image" && m.URL).map((m) => m.URL)
    : [];

  const handleLikeOnClicked = async () => {
    // POST like api
    if (!isLiked) {
      const res = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}${process.env.NEXT_PUBLIC_POSTS_API}/${post_id}/like`,
        {
          method: "PUT",
          credentials: "include",
        }
      );
      if (!res.ok) {
        console.error("Failed to like the post");
        return;
      }
      setLikes((prevLikes) =>
        prevLikes === like_count ? prevLikes + 1 : like_count
      );
      setIsLiked(true); // 更新 isLiked 狀態
    } else {
      const res = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}${process.env.NEXT_PUBLIC_POSTS_API}/${post_id}/unlike`,
        {
          method: "PUT",
          credentials: "include",
        }
      );
      if (!res.ok) {
        console.error("Failed to unlike the post");
        return;
      }
      setLikes((prevLikes) =>
        prevLikes === like_count ? prevLikes - 1 : like_count
      );
      setIsLiked(false); // 更新 isLiked 狀態
    }
  };

  const handleComment = () => {
    // 可根據需求實作
  };

  return (
    <article className="bg-white p-6 w-full border-t border-gray-200">
      <div className="flex items-center p-2">
        <img
          src={avatar_url}
          alt={`${author_name}'s avatar`}
          className="w-12 h-12 rounded-full mr-4 object-cover ring-1 ring-offset-1 ring-[#B6B09F]"
        />
        <div>
          <Link href={`/profile/${post.author_id}`}>
            <h4 className="font-semibold text-gray-800">{author_name}</h4>
          </Link>

          <p className="text-sm text-gray-500">
            發布於：
            {typeof created_at === "string"
              ? created_at
              : created_at?.toLocaleString?.() || ""}
          </p>
        </div>
      </div>
      <div className="prose prose-indigo max-w-none p-2 mb-4">
        <p>{content}</p>
      </div>
      {/* 圖片展示區塊 */}
      {imageSrcs.length > 0 && (
        <div className="overflow-x-auto flex flex-row gap-2 p-2">
          {imageSrcs.map((src) => (
            <Image
              key={src}
              src={src}
              width={240}
              height={160}
              className="rounded-lg object-cover"
              alt="Post media"
            />
          ))}
        </div>
      )}
      <div className="mt-4 pt-4 border-t border-gray-200">
        <div className="flex items-center space-x-4">
          <button
            onClick={handleLikeOnClicked}
            className="flex items-center text-gray-600 hover:text-[#B6B09F] focus:outline-none transition-colors duration-150"
            aria-label={`Like this post. Currently ${likes} likes.`}
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill={isLiked ? "currentColor" : "none"}
              viewBox="0 0 24 24"
              strokeWidth={1.5}
              stroke="currentColor"
              className={`size-6 mr-1.5 ${isLiked ? "text-red-500" : ""}`}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12Z"
              />
            </svg>
            <span className="font-medium">{likes}</span>
            <span className="ml-1 hidden sm:inline">Likes</span>
          </button>
          <button
            onClick={handleComment}
            className="flex items-center text-gray-600 hover:text-[#B6B09F] transition-colors duration-150"
            aria-label={`View comments. Currently ${comment_count} comments.`}
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              strokeWidth={1.5}
              stroke="currentColor"
              className="size-6 mr-1.5"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M7.5 8.25h9m-9 3H12m-9.75 1.51c0 1.6 1.123 2.994 2.707 3.227 1.129.166 2.27.293 3.423.379.35.026.67.21.865.501L12 21l2.755-4.133a1.14 1.14 0 0 1 .865-.501 48.172 48.172 0 0 0 3.423-.379c1.584-.233 2.707-1.626 2.707-3.228V6.741c0-1.602-1.123-2.995-2.707-3.228A48.394 48.394 0 0 0 12 3c-2.392 0-4.744.175-7.043.513C3.373 3.746 2.25 5.14 2.25 6.741v6.018Z"
              />
            </svg>
            <span className="font-medium">{comment_count}</span>
            <span className="ml-1 hidden sm:inline">Comments</span>
          </button>
        </div>
      </div>
    </article>
  );
}
