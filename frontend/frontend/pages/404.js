import Image from 'next/image';
import { ConstructingNavbar } from '@/components/navbar';
import Layout from '@/components/layout';
import { useRouter } from 'next/router';
import { useState, useEffect } from 'react';

export default function Constructing() {
    const router = useRouter();
    const [countdown, setCountdown] = useState(5); // 10 seconds countdown

    useEffect(() => {
        const timer = setInterval(() => {
            setCountdown((prev) => {
                if (prev <= 1) {
                    router.push('/');
                    return 0;
                }
                return prev - 1;
            });
        }, 1000);

        return () => clearInterval(timer);
    }, [router]);

    return (
        <Layout pageTitle="404 Not Found">
            <ConstructingNavbar />
            <div className="flex flex-col items-center justify-center min-h-screen bg-[#F2F2F2]">
                <h1 className="text-4xl font-bold mb-4">正在建設中</h1>
                <p className="text-lg text-gray-700 mb-4">這個頁面正在建設中，敬請期待！</p>
                <p className="text-md text-gray-600 mb-8">
                    將在 {countdown} 秒後重新導向至首頁
                </p>
                <Image
                    src="/constructing.png"
                    alt="Under Construction"
                    width={256}
                    height={256}
                    className="w-64 h-auto"
                />
            </div>
        </Layout>
    );
}