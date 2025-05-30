import Image from "next/image";
import Link from 'next/link';

export function NavigationCard({ title, description, link, bgColor = "bg-indigo-500" }) {
  return (
    <Link href={link} legacyBehavior>
      <a className={`block p-6 rounded-lg shadow-lg text-white hover:opacity-90 transition-opacity ${bgColor}`}>
        <h3 className="text-xl font-bold mb-2">{title}</h3>
        <p className="text-sm">{description}</p>
      </a>
    </Link>
  );
}

export function PostCard(props) {
  let postedBy = props.postedBy ? props.postedBy : "Default User";
  let postedByAvatar = props.postedByAvatar ? props.postedByAvatar : "https://images.unsplash.com/photo-1491528323818-fdd1faba62cc?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80";
  let postedByLink = props.postedByLink ? props.postedByLink : "#";

  let postedAt = props.postedAt ? props.postedAt : "2023-10-01 12:00:00";
  let contents = props.contents ? props.contents : "這是貼文的主要內容。Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur et leveraging synergies to enhance stakeholder engagement.";
  let imageSrcs = props.srcs ? props.srcs : ["/finger.png"];
  let likes = props.likes ? props.likes : 14;
  let comments = props.comments ? props.comments : ["這是第一條評論。", "這是第二條評論。"];
  let numberOfComments = comments.length;

  return (
    <main className="lg:col-span-6 md:col-span-8 col-span-12">
      {" "}
      {/* 確保在不同斷點下有合理的寬度 */}
      <article className="bg-white p-6 rounded-lg shadow-lg">
        <div className="flex items-center mb-4">
          <img
            src={postedByAvatar}
            alt={postedBy}
            className="w-12 h-12 rounded-full mr-4"
          />
          <div>
            <h4 className="font-semibold text-gray-800"></h4>
            <p className="text-sm text-gray-500">發布於：{postedAt}</p>
          </div>
        </div>
        {/* 使用 prose 類別可以快速美化文章內容排版 */}
        <div className="prose prose-indigo max-w-none">
          <p>
            {contents}
          </p>
          <br />
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

          </div>
          <div className="mt-6 pt-4 border-t border-gray-200">
            {/* 互動按鈕，例如按讚、留言 */}
            <div className="flex items-center space-x-6"> {/* Refactored: main flex container for buttons with spacing */}
            {/* Likes button group */}
            <button className="flex items-center text-[#0D0D0D] hover:text-gray-700 focus:outline-none">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6 mr-2">
              <path strokeLinecap="round" strokeLinejoin="round" d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12Z" />
              </svg>
              <span>Likes ({likes})</span>
            </button>

            {/* Comments button group */}
            <button className="flex items-center text-[#0D0D0D] hover:text-gray-700 focus:outline-none">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6 mr-2">
              <path strokeLinecap="round" strokeLinejoin="round" d="M7.5 8.25h9m-9 3H12m-9.75 1.51c0 1.6 1.123 2.994 2.707 3.227 1.129.166 2.27.293 3.423.379.35.026.67.21.865.501L12 21l2.755-4.133a1.14 1.14 0 0 1 .865-.501 48.172 48.172 0 0 0 3.423-.379c1.584-.233 2.707-1.626 2.707-3.228V6.741c0-1.602-1.123-2.995-2.707-3.228A48.394 48.394 0 0 0 12 3c-2.392 0-4.744.175-7.043.513C3.373 3.746 2.25 5.14 2.25 6.741v6.018Z" />
              </svg>
              <span>Comments ({numberOfComments})</span>
            </button>
            </div>
          </div>
          </article>
          {/* 如果有多個貼文，可以在這裡繼續添加 */}
    </main>
  );
}
