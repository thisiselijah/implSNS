// import Image from "next/image";
import Layout from "@/components/layout";
import { IndexNavbar } from "@/components/navbar";
import { Geist, Geist_Mono } from "next/font/google";
import {NavigationCard} from "@/components/card"; // 引入我們剛才建立的卡片組件

import { Swiper, SwiperSlide } from 'swiper/react';
import 'swiper/css'; 
import { Autoplay, Pagination, Navigation } from 'swiper/modules';
import 'swiper/css/pagination';
import 'swiper/css/navigation';

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

// 導覽卡片的資料
const navCardsData = [
  { id: 1, title: "探索功能", description: "發現我們所有的強大功能。", link: "/constructing", bgColor: "bg-sky-500" },
  { id: 2, title: "了解趨勢", description: "看看大家都在流行什麼？", link: "/constructing", bgColor: "bg-emerald-500" },
  { id: 3, title: "關於我們", description: "了解更多我們的團隊。", link: "/constructing", bgColor: "bg-purple-500" },
  { id: 4, title: "聯絡客服", description: "需要協助嗎？隨時與我們聯繫。", link: "/constructing", bgColor: "bg-amber-500" },
];

export default function Home() {
  return (
    <Layout pageTitle="首頁"> {/* 假設 Layout 組件可以接受 pageTitle */}
      <IndexNavbar />
      
      {/* 自動翻頁的導覽卡片區塊 */}
      <section className={`py-12 sm:py-16 ${geistSans.variable} ${geistMono.variable}`}>
        <div className="container mx-auto px-4">
          <h2 className="text-3xl font-bold text-center text-gray-800 mb-10">快速導覽</h2>
          
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
              640: { // sm
                slidesPerView: 2,
                spaceBetween: 20,
              },
              768: { // md
                slidesPerView: 2,
                spaceBetween: 30,
              },
              1024: { // lg
                slidesPerView: 3,
                spaceBetween: 30,
              },
              1280: { // xl
                slidesPerView: 4,
                spaceBetween: 30,
              }
            }}
          >
            {navCardsData.map(card => (
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
      </section>

      {/* 你原本的 columns-2xs div，如果不需要可以移除 */}
      {/* <div className="columns-2xs"></div> */}

      {/* 其他首頁內容可以放在這裡 */}
      <div className="container mx-auto px-4 py-16 text-center">
        <h1 className="text-4xl font-bold">歡迎來到我的網站</h1>
        <p className="mt-4 text-lg text-gray-600">資工三 S11159050</p>
      </div>
    </Layout>
  );
}
