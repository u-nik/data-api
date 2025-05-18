'use client';

import {useEffect} from 'react';
import {useAuth} from '@/components/auth/AuthContext';
import {usePathname, useRouter} from 'next/navigation';

export default function AuthGuard({children}: {children: React.ReactNode}) {
    const {accessToken} = useAuth();
    const router = useRouter();
    const pathname = usePathname();

    useEffect(() => {
        if (!accessToken) {
            // Save current path for redirect after login
            sessionStorage.setItem('post_login_redirect', pathname);
            router.replace('/login');
        }
    }, [accessToken, pathname, router]);

    if (!accessToken) return null;
    return <>{children}</>;
}
