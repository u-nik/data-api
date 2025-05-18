import React from 'react';
import {AuthLayout} from '@/components/ui/auth-layout';
import {Logo} from '@/components/ui/logo';
import {AuthProvider} from '@/components/auth/AuthContext';

export default function RootLayout({
    children,
}: Readonly<{children: React.ReactNode}>) {
    return (
        <AuthProvider>
            <AuthLayout>
                <div className='grid w-full max-w-sm grid-cols-1 gap-8'>
                    <Logo src='/logo.png' />
                    {children}
                </div>
            </AuthLayout>
        </AuthProvider>
    );
}
