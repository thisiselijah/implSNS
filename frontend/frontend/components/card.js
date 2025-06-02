import Image from "next/image";
import Link from "next/link"; // Link 用於 NavigationCard
import { useState } from "react";

// NavigationCard 組件保持不變，但建議移除 legacyBehavior (參考之前的建議)
export function NavigationCard({
  title,
  description,
  link,
  bgColor = "bg-indigo-500",
}) {
  return (
    // 建議：移除 legacyBehavior 並將 className 移至 Link
    <Link
      href={link}
      className={`block p-6 rounded-lg shadow-lg text-white hover:opacity-90 transition-opacity ${bgColor}`}
    >
      <h3 className="text-xl font-bold mb-2">{title}</h3>
      <p className="text-sm">{description}</p>
    </Link>
  );
}

export function PostCard(props) {
  // 從 props 或預設值初始化狀態
  const initialLikes = props.likes !== undefined ? props.likes : 14;
  const initialComments =
    props.comments !== undefined
      ? props.comments
      : ["這是第一條評論。", "這是第二條評論。"];

  const [likes, setLikes] = useState(initialLikes);
  // 假設 comments 狀態將用於顯示或添加評論
  const [commentsData, setCommentsData] = useState(initialComments);

  const postedBy = props.postedBy || "Default User";
  const postedByAvatar =
    props.postedByAvatar ||
    "https://images.unsplash.com/photo-1491528323818-fdd1faba62cc?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80";
  // postedByLink 未在目前模板中使用，但保留
  // const postedByLink = props.postedByLink || "#";
  const postedAt = props.postedAt || "2023-10-01 12:00:00";
  const contents =
    props.contents ||
    "這是貼文的主要內容。Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur et leveraging synergies to enhance stakeholder engagement.";
  const imageSrcs = props.srcs || ["/finger.png"]; // 確保路徑正確，例如 '/images/finger.png'
  const numberOfComments = commentsData.length;

  const handleLike = () => {
    setLikes((prevLikes) => 
      initialLikes === prevLikes ? prevLikes + 1 : prevLikes - 1); // 如果 likes 大於初始值，則增加 likes
 
  };

  const handleComment = () => {
    // 這裡可以添加打開評論框或跳轉到評論區的邏輯
    console.log(
      "Comment button clicked. Number of comments:",
      numberOfComments
    );
    // 示例：添加一條新評論 (實際應用中會有輸入框)
    // setCommentsData(prevComments => [...prevComments, "新評論!"]);
  };

  return (
    // PostCard 不應該自己決定 grid span，這應該由父組件 (如 post.js) 在 grid 佈局中指定
    // <main className="lg:col-span-6 md:col-span-8 col-span-12">
    <article className="bg-white p-6 rounded-lg shadow-lg w-full">
      {" "}
      {/* w-full 使其填滿父容器指定的寬度 */}
      <div className="flex items-center mb-4">
        <img
          src={postedByAvatar}
          alt={`${postedBy}'s avatar`}
          className="w-12 h-12 rounded-full mr-4 object-cover" // object-cover 避免圖片變形
        />
        <div>
          {/* 你可能想在這裡顯示 postedBy */}
          <h4 className="font-semibold text-gray-800">{postedBy}</h4>
          <p className="text-sm text-gray-500">發布於：{postedAt}</p>
        </div>
      </div>
      <div className="prose prose-indigo max-w-none mb-4">
        {" "}
        {/* mb-4 給內容和圖片區塊一些底部間距 */}
        <p>{contents}</p>
      </div>
      {/* 圖片展示區塊 */}
      <div>
        {imageSrcs.map((src) => {
          return (
            <Image
              key={src}
              src={src}
              alt="貼文圖片"
              width={300}
              height={200}
              className="rounded-lg mb-4"
            />
          );
        })}
      </div>
      <div className="mt-6 pt-4 border-t border-gray-200">
        <div className="flex items-center space-x-4">
          {" "}
          {/* 調整 space-x */}
          <button
            onClick={handleLike}
            className="flex items-center text-gray-600 hover:text-[#B6B09F] focus:outline-none transition-colors duration-150"
            aria-label={`Like this post. Currently ${likes} likes.`}
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill={likes > initialLikes ? "currentColor" : "none"}
              viewBox="0 0 24 24"
              strokeWidth={1.5}
              stroke="currentColor"
              className={`size-6 mr-1.5 ${
                likes > initialLikes ? "text-red-500" : ""
              }`}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12Z"
              />
            </svg>
            <span className="font-medium">{likes}</span>
            <span className="ml-1 hidden sm:inline">Likes</span>{" "}
            {/* 在小螢幕上可選隱藏 "Likes" 文字 */}
          </button>
          {/* Comments 按鈕 */}
          <button
            onClick={handleComment}
            className="flex items-center text-gray-600 hover:text-[#B6B09F] transition-colors duration-150"
            aria-label={`View comments. Currently ${numberOfComments} comments.`}
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
            <span className="font-medium">{numberOfComments}</span>
            <span className="ml-1 hidden sm:inline">Comments</span>{" "}
            {/* 在小螢幕上可選隱藏 "Comments" 文字 */}
          </button>
          {/* 你可以在這裡添加更多互動按鈕，例如分享 */}
        </div>
      </div>
    </article>
    // </main>
  );
}
