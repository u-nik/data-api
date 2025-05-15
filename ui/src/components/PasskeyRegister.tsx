'use client';

import {useState} from 'react';
import {webauthnRegister} from './WebAuthn';

export default function PasskeyRegister() {
    const [username, setUsername] = useState('');
    const [message, setMessage] = useState('');
    const [registered, setRegistered] = useState(false);

    const handleRegister = async () => {
        const ok = await webauthnRegister(username, setMessage);
        if (ok) setRegistered(true);
    };

    if (registered) {
        return (
            <div>
                Registration successful! You can now{' '}
                <a href='/login' className='text-blue-600 underline'>
                    login
                </a>
                .
            </div>
        );
    }

    return (
        <div className='max-w-md mx-auto mt-12 p-8 bg-white rounded-xl shadow-md border border-gray-200'>
            <h2 className='text-2xl font-bold mb-6 text-center'>
                Register with Passkey
            </h2>
            <input
                type='text'
                placeholder='Username'
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                className='w-full mb-4 px-4 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-400'
            />
            <div className='flex gap-4 mb-4'>
                <button
                    onClick={handleRegister}
                    className='flex-1 bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded transition-colors'
                >
                    Register
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
