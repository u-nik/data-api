import React from 'react';
import {AuthLayout} from '@/components/catalyst/auth-layout';

export default function RootLayout({
    children,
}: Readonly<{children: React.ReactNode}>) {
    return <AuthLayout>{children}</AuthLayout>;
}
