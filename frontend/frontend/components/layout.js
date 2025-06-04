import Head from "next/head";
import styles from "./layout.module.css"; // 確保這個路徑正確，指向你的 CSS 模組檔案


// import styles from "./layout.module.css"; 

export default function Layout({ children, pageTitle}) { 
    return (
        <div className={styles.layout}> {/* 1. 設定為 flex 容器，垂直排列，最小高度為螢幕高度，加上一個淺灰色背景 */}
            <Head>
                <link rel="icon" href="/favicon.ico" />
                <meta name="description" content="Social media project for cloud computing course" />
                <meta name="viewport" content="width=device-width, initial-scale=1" />
                {/* 2. 修改 title 標籤的用法，使其能動態載入或使用預設值 */}
                <title>{pageTitle}</title>
            </Head>

            {/* <header className="bg-white shadow-sm p-4">
                <div className="container mx-auto">
                    <p className="text-xl font-semibold">My App Header</p>
                </div>
            </header> */}

            {/* 4. 設定 main 區域可以成長以填滿剩餘空間，並加上一些基本內邊距和容器設定 */}
            <main className="flex-grow">
                {/*  container mx-auto p-4 md:p-6 lg:p-8 */}
                {children}
            </main>

            {/* 5. 為 footer 加上樣式，並確保文字居中 */}
            <footer className="p-4 text-center mt-auto ">
                <p className="text-sm">&copy; 2025 Copywrite Social Media Project. All rights reserved.</p>
                
            </footer>
        </div>
    );
}