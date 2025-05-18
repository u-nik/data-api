'use client';

import {useState, useEffect} from 'react';
import {useAuth} from '@/components/auth/AuthContext';
import {useRouter} from 'next/navigation';
import {Button} from '@/components/ui/button';

export default function LoginPage() {
    const {login, accessToken} = useAuth();
    const router = useRouter();
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const loginError = sessionStorage.getItem('login_error');
        if (loginError) {
            setError(loginError);
            sessionStorage.removeItem('login_error');
        }
    }, []);

    useEffect(() => {
        if (accessToken) {
            const redirect =
                sessionStorage.getItem('post_login_redirect') || '/';
            sessionStorage.removeItem('post_login_redirect');
            router.replace(redirect);
        }
    }, [accessToken, router]);

    return (
        <div className='flex flex-col items-center justify-center'>
            {error && <div className='mb-4 text-red-600'>{error}</div>}
            <Button onClick={login}>Anmelden</Button>
        </div>
    );
}
