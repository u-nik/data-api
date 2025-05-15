import type {NextConfig} from 'next';

const nextConfig: NextConfig = {
    async rewrites() {
        return [
            {
                source: '/api/:path*',
                destination: 'http://localhost:8080/api/:path*', // Passe ggf. den Port an dein Go-Backend an
            },
        ];
    },
};

export default nextConfig;
