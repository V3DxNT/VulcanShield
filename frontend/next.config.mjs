/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  async rewrites() {
    return [
      {
        source: '/api/v1/:path*',
        destination: 'http://backend:8080/api/v1/:path*',
      },
    ];
  },
};

export default nextConfig;
