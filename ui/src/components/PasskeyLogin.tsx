import React, {useState} from 'react';
import {startAuthentication, startRegistration} from '@simplewebauthn/browser';
import axios from 'axios';

export default function PasskeyLogin() {
    const [username, setUsername] = useState('');
    const [message, setMessage] = useState('');
    const [isLoggedIn, setIsLoggedIn] = useState(false);

    // Registration handler
    const handleRegister = async () => {
        setMessage('');
        try {
            // 1. Get registration options from backend
            const {data: options} = await axios.post(
                '/api/webauthn/register/options',
                {username},
            );
            // 2. Start registration ceremony
            const attResp = await startRegistration(options);
            // 3. Send attestation response to backend for verification
            await axios.post('/api/webauthn/register/verify', {
                username,
                attestation: attResp,
            });
            setMessage('Registration successful! You can now log in.');
        } catch (err: unknown) {
            if (axios.isAxiosError(err)) {
                setMessage('Registration failed: ' + (err?.message || err));
            }
        }
    };

    // Login handler
    const handleLogin = async () => {
        setMessage('');
        try {
            // 1. Get authentication options from backend
            const {data: options} = await axios.post(
                '/api/webauthn/login/options',
                {username},
            );
            // 2. Start authentication ceremony
            const assertionResp = await startAuthentication(options);
            // 3. Send assertion response to backend for verification
            await axios.post('/api/webauthn/login/verify', {
                username,
                assertion: assertionResp,
            });
            setIsLoggedIn(true);
            setMessage('Login successful!');
        } catch (err: unknown) {
            if (axios.isAxiosError(err)) {
                setMessage('Login failed: ' + (err?.message || err));
            }
        }
    };

    if (isLoggedIn) {
        return (
            <div>Welcome, {username}! You are logged in with a Passkey.</div>
        );
    }

    return (
        <div
            style={{
                maxWidth: 400,
                margin: '2rem auto',
                padding: 24,
                border: '1px solid #ccc',
                borderRadius: 8,
            }}
        >
            <h2>Login/Register with Passkey</h2>
            <input
                type='text'
                placeholder='Username'
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                style={{width: '100%', marginBottom: 12, padding: 8}}
            />
            <div style={{display: 'flex', gap: 8}}>
                <button onClick={handleRegister}>Register</button>
                <button onClick={handleLogin}>Login</button>
            </div>
            {message && (
                <div
                    style={{
                        marginTop: 16,
                        color: message.includes('failed') ? 'red' : 'green',
                    }}
                >
                    {message}
                </div>
            )}
        </div>
    );
}
