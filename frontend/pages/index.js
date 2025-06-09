// import Image from "next/image";
import Layout from "@/components/layout";
import { IndexNavbar } from "@/components/navbar";
// import { Geist, Geist_Mono } from "next/font/google";
import { NavigationCard } from "@/components/card";

// import { Swiper, SwiperSlide } from "swiper/react";
// import "swiper/css";
// import { Autoplay, Pagination, Navigation } from "swiper/modules";
// import "swiper/css/pagination";
// import "swiper/css/navigation";
import { useEffect, useRef } from "react";
import { animate, createScope } from "animejs";

// const geistSans = Geist({
//   variable: "--font-geist-sans",
//   subsets: ["latin"],
// });

// const geistMono = Geist_Mono({
//   variable: "--font-geist-mono",
//   subsets: ["latin"],
// });

// // 導覽卡片的資料
// const navCardsData = [
//   {
//     id: 1,
//     title: "探索功能",
//     description: "發現我們所有的強大功能。",
//     link: "/constructing",
//     bgColor: "bg-[#B6B09F]",
//   },
//   {
//     id: 2,
//     title: "了解趨勢",
//     description: "看看大家都在流行什麼？",
//     link: "/constructing",
//     bgColor: "bg-[#B6B09F]",
//   },
//   {
//     id: 3,
//     title: "關於我們",
//     description: "了解更多我們的團隊。",
//     link: "/constructing",
//     bgColor: "bg-[#B6B09F]",
//   },
//   {
//     id: 4,
//     title: "聯絡客服",
//     description: "需要協助嗎？隨時與我們聯繫。",
//     link: "/constructing",
//     bgColor: "bg-[#B6B09F]",
//   },
// ];

export default function Home() {
  // 使用 useRef 來獲取 h2 元素
  const spanRef = useRef(null);
  const scope = useRef(null);

  useEffect(() => {
    scope.current = createScope({ spanRef }).add((self) => {
      animate("span", {
        // Property keyframes
        y: [
          { to: "-2.75rem", ease: "outExpo", duration: 600 },
          { to: 0, ease: "outBounce", duration: 800, delay: 100 },
        ],
        // Property specific parameters
        rotate: {
          from: "-1turn",
          delay: 0,
        },
        delay: (_, i) => i * 50, // Function based value
        ease: "inOutCirc",
        loopDelay: 1000,
        loop: true,
      });
    });
    return () => scope.current.revert();
  }, []);

  return (
    <Layout pageTitle="首頁">
      <IndexNavbar />

      {/* 自動翻頁的導覽卡片區塊 */}
      {/* <section
        className={`py-12 sm:py-20 ${geistSans.variable} ${geistMono.variable}`}
      >
        <div className="container mx-auto px-10">
          <h2 className="text-3xl font-bold text-center text-gray-800 mb-10">
            快速導覽
          </h2>

          <Swiper
            // 3. 註冊需要的模組
            modules={[Autoplay, Pagination, Navigation]}
            spaceBetween={30} // Slide 之間的間距
            slidesPerView={3} // 預設顯示一個 Slide
            loop={true} // 開啟無限循環
            autoplay={{
              delay: 3000, // 自動播放的延遲時間 (毫秒)
              disableOnInteraction: false, // 使用者互動後是否停止自動播放 (false 表示不停止)
              pauseOnMouseEnter: true, // 滑鼠移入時暫停自動播放
            }}
            pagination={{
              clickable: true, // 分頁點可以點擊切換
            }}
            navigation={false} // 開啟上一個/下一個導覽箭頭 (可選)
            className="mySwiper " // 你可以添加自訂 class 來進一步調整樣式
            // 響應式設定：不同螢幕尺寸下顯示不同數量的 Slide
            breakpoints={{
              640: {
                // sm
                slidesPerView: 2,
                spaceBetween: 20,
              },
              768: {
                // md
                slidesPerView: 2,
                spaceBetween: 30,
              },
              1024: {
                // lg
                slidesPerView: 3,
                spaceBetween: 30,
              },
              1280: {
                // xl
                slidesPerView: 4,
                spaceBetween: 30,
              },
            }}
          >
            {navCardsData.map((card) => (
              <SwiperSlide key={card.id}>
                <NavigationCard
                  title={card.title}
                  description={card.description}
                  link={card.link}
                  bgColor={card.bgColor}
                />
              </SwiperSlide>
            ))}
          </Swiper>
        </div>
      </section> */}

      <div className="flex items-center justify-center  min-h-screen">
        <div className="text-center">
          <h1 className="text-5xl grid grid-cols-8 gap-2 place-items-center">
            {["歡", "迎", "來", "到", "我", "的", "網", "站"].map(
              (char, idx) => (
                <span key={idx}>{char}</span>
              )
            )}
          </h1>
          <p className="mt-4 text-lg text-gray-600">資工三 S11159050</p>
        </div>
      </div>
    </Layout>
  );
}
