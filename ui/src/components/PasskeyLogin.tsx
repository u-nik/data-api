'use client';

import {useEffect, useRef, useState} from 'react';
import {webauthnLogin} from './WebAuthn';

export default function PasskeyLogin() {
    const [username, setUsername] = useState('');
    const [message, setMessage] = useState('');
    const [isLoggedIn, setIsLoggedIn] = useState(false);
    const [conditional, setConditional] = useState(true);
    const inputRef = useRef<HTMLInputElement>(null);

    // Autofill: trigger passkey autofill on mount if conditional mediation is enabled
    useEffect(() => {
        if (conditional) {
            handleLogin();
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [conditional]);

    // Modified login handler for conditional UI
    const handleLogin = async () => {
        // If conditional mediation, do not require username
        const uname = conditional ? '' : username;
        await webauthnLogin(
            uname,
            setMessage,
            (v) => {
                setIsLoggedIn(v);
                // Optionally, set username from backend response if available
            },
            conditional,
        );
        // Optionally, set username from backend response if available
        // (requires backend to return username)
    };

    if (isLoggedIn) {
        return (
            <div>
                Welcome{username ? `, ${username}` : ''}! You are logged in with
                a Passkey.
            </div>
        );
    }

    return (
        <div className='max-w-md mx-auto mt-12 p-8 bg-white rounded-xl shadow-md border border-gray-200'>
            <h2 className='text-2xl font-bold mb-6 text-center'>
                Login with Passkey
            </h2>
            <input
                ref={inputRef}
                type='text'
                placeholder='Username (not required for passkey autofill)'
                value={username}
                autoComplete='username webauthn'
                onChange={(e) => setUsername(e.target.value)}
                className='w-full mb-4 px-4 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-400'
                disabled={conditional}
            />
            <div className='flex items-center mb-4'>
                <input
                    id='conditional-mediation'
                    type='checkbox'
                    checked={conditional}
                    onChange={() => setConditional((v) => !v)}
                    className='mr-2 accent-blue-500'
                />
                <label htmlFor='conditional-mediation' className='text-sm'>
                    Use passkey autofill (conditional mediation)
                </label>
            </div>
            <div className='flex gap-4 mb-4'>
                <button
                    onClick={handleLogin}
                    className='flex-1 bg-green-500 hover:bg-green-600 text-white font-semibold py-2 px-4 rounded transition-colors'
                    disabled={conditional} // Only allow manual login if not conditional
                >
                    Login
                </button>
            </div>
            {message && (
                <div
                    className={`mt-4 text-center font-medium ${
                        message.includes('failed')
                            ? 'text-red-600'
                            : 'text-green-600'
                    }`}
                >
                    {message}
                </div>
            )}
        </div>
    );
}
