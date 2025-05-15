import {startAuthentication, startRegistration} from '@simplewebauthn/browser';
import axios from 'axios';

export async function webauthnRegister(
    username: string,
    setMessage: (msg: string) => void,
) {
    setMessage('');
    try {
        const {data: options} = await axios.post(
            '/api/webauthn/register/options',
            {username},
        );
        const attResp = await startRegistration(options);
        await axios.post('/api/webauthn/register/verify', {
            username,
            attestation: attResp,
        });
        setMessage('Registration successful! You can now log in.');
        return true;
    } catch (err: unknown) {
        if (axios.isAxiosError(err)) {
            setMessage('Registration failed: ' + (err?.message || err));
        }
        return false;
    }
}

export async function webauthnLogin(
    username: string,
    setMessage: (msg: string) => void,
    setIsLoggedIn: (v: boolean) => void,
) {
    setMessage('');
    try {
        const {data: options} = await axios.post(
            '/api/webauthn/login/options',
            {username},
        );
        const assertionResp = await startAuthentication(options);
        await axios.post('/api/webauthn/login/verify', {
            username,
            assertion: assertionResp,
        });
        setIsLoggedIn(true);
        setMessage('Login successful!');
        return true;
    } catch (err: unknown) {
        if (axios.isAxiosError(err)) {
            setMessage('Login failed: ' + (err?.message || err));
        }
        return false;
    }
}
