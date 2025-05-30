// import Image from "next/image";
import Layout from "@/components/layout";
import Navbar from "@/components/navbar";
import { Geist, Geist_Mono } from "next/font/google";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export default function Home() {
  return (
    <Layout>
      <Navbar />
      <div className="columns-2xs gap-4">
        <p>Hello World!</p>
        <p>妳好</p>
      </div>
    </Layout>
  );
}
