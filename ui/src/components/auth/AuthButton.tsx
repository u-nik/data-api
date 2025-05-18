import React from 'react';
import {useAuth} from './AuthContext';

export const AuthButton: React.FC = () => {
    const {accessToken, login, logout} = useAuth();
    if (accessToken) {
        return <button onClick={logout}>Logout</button>;
    }
    return <button onClick={login}>Login</button>;
};
