'use client';

import {useEffect} from 'react';
import {useRouter} from 'next/navigation';
import {useAuth} from '@/components/auth/AuthContext';

const OAUTH2_TOKEN_URL = process.env.NEXT_PUBLIC_OAUTH2_TOKEN_URL!;
const OAUTH2_CLIENT_ID = process.env.NEXT_PUBLIC_OAUTH2_CLIENT_ID!;
const OAUTH2_REDIRECT_URI = process.env.NEXT_PUBLIC_OAUTH2_REDIRECT_URI!;

export default function AuthCallback() {
    const router = useRouter();
    const {accessToken, login} = useAuth();

    useEffect(() => {
        const url = new URL(window.location.href);
        const code = url.searchParams.get('code');
        const error = url.searchParams.get('error');
        const errorDescription = url.searchParams.get('error_description');
        if (error) {
            console.error('OAuth2 Error:', error, errorDescription);
            sessionStorage.setItem(
                'login_error',
                'Beim Login trat ein Fehler auf.',
            );
            router.replace('/login');
            return;
        }
        const codeVerifier = sessionStorage.getItem('pkce_code_verifier');
        if (!code || !codeVerifier) return;

        async function exchangeCode() {
            try {
                const params = new URLSearchParams({
                    grant_type: 'authorization_code',
                    code: code as string,
                    redirect_uri: OAUTH2_REDIRECT_URI,
                    client_id: OAUTH2_CLIENT_ID,
                    code_verifier: codeVerifier as string,
                });
                const res = await fetch(OAUTH2_TOKEN_URL, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded',
                    },
                    body: params.toString(),
                });
                if (res.ok) {
                    const data = await res.json();
                    window.dispatchEvent(
                        new CustomEvent('auth-token', {
                            detail: data.access_token,
                        }),
                    );
                    router.replace('/');
                } else {
                    console.error('Token exchange failed (2):', res.statusText);
                    sessionStorage.setItem(
                        'login_error',
                        'Beim Login trat ein Fehler auf.',
                    );
                    router.replace('/login');
                }
            } catch {
                console.error('Token exchange failed (1):', error);
                sessionStorage.setItem(
                    'login_error',
                    'Beim Login trat ein Fehler auf.',
                );
                router.replace('/login');
            }
        }
        exchangeCode();
    }, [router, login, accessToken]);

    return <div>Authenticating...</div>;
}
