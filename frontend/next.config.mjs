/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  reactStrictMode: true,
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: 'media-bucket-thisiselijah-implsns.s3.us-east-1.amazonaws.com',
        pathname: '/**',
      },
    ],
  },
};

export default nextConfig;