'use client';

import React, {
    createContext,
    useContext,
    useState,
    useCallback,
    useEffect,
} from 'react';

interface AuthContextType {
    accessToken: string | null;
    login: () => void;
    logout: () => void;
}

const AuthContext = createContext<AuthContextType>({
    accessToken: null,
    login: () => {},
    logout: () => {},
});

export const useAuth = () => useContext(AuthContext);

const OAUTH2_CLIENT_ID = process.env.NEXT_PUBLIC_OAUTH2_CLIENT_ID || '';
const OAUTH2_AUTH_URL = process.env.NEXT_PUBLIC_OAUTH2_AUTH_URL || '';
const OAUTH2_REDIRECT_URI = process.env.NEXT_PUBLIC_OAUTH2_REDIRECT_URI || '';
const OAUTH2_SCOPE = 'openid profile email';

function generateCodeVerifier(length = 128) {
    const array = new Uint8Array(length);
    window.crypto.getRandomValues(array);
    return btoa(String.fromCharCode(...array))
        .replace(/\+/g, '-')
        .replace(/\//g, '_')
        .replace(/=+$/, '');
}

async function generateCodeChallenge(verifier: string) {
    const encoder = new TextEncoder();
    const data = encoder.encode(verifier);
    const digest = await window.crypto.subtle.digest('SHA-256', data);
    return btoa(String.fromCharCode(...new Uint8Array(digest)))
        .replace(/\+/g, '-')
        .replace(/\//g, '_')
        .replace(/=+$/, '');
}

export const AuthProvider: React.FC<{children: React.ReactNode}> = ({
    children,
}) => {
    const [accessToken, setAccessToken] = useState<string | null>(null);

    // Listen for token from callback page
    useEffect(() => {
        const handler = (e: CustomEvent<string>) => {
            if (e.detail) setAccessToken(e.detail);
        };
        window.addEventListener('auth-token', handler as EventListener);
        return () =>
            window.removeEventListener('auth-token', handler as EventListener);
    }, []);

    const login = useCallback(async () => {
        console.log('login clicked');
        console.log('OAUTH2_CLIENT_ID', OAUTH2_CLIENT_ID);
        console.log('OAUTH2_AUTH_URL', OAUTH2_AUTH_URL);
        console.log('OAUTH2_REDIRECT_URI', OAUTH2_REDIRECT_URI);
        const codeVerifier = generateCodeVerifier();
        const codeChallenge = await generateCodeChallenge(codeVerifier);
        sessionStorage.setItem('pkce_code_verifier', codeVerifier);
        // Generate a random state with at least 8 characters
        const state = btoa(
            String.fromCharCode(
                ...window.crypto.getRandomValues(new Uint8Array(16)),
            ),
        )
            .replace(/\+/g, '-')
            .replace(/\//g, '_')
            .replace(/=+$/, '')
            .slice(0, 16);
        sessionStorage.setItem('oauth2_state', state);
        const params = new URLSearchParams({
            client_id: OAUTH2_CLIENT_ID,
            response_type: 'code',
            redirect_uri: OAUTH2_REDIRECT_URI,
            scope: OAUTH2_SCOPE,
            code_challenge: codeChallenge,
            code_challenge_method: 'S256',
            state,
        });
        const url = `${OAUTH2_AUTH_URL}?${params.toString()}`;
        console.log('Login redirect URL:', url);
        window.location.href = url;
    }, []);

    const logout = useCallback(() => {
        setAccessToken(null);
    }, []);

    return (
        <AuthContext.Provider value={{accessToken, login, logout}}>
            {children}
        </AuthContext.Provider>
    );
};
