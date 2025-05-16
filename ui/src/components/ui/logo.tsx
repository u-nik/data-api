import React from 'react';
import Image from 'next/image';

export function Logo({...props}: React.ComponentPropsWithoutRef<'img'>) {
    return (
        <Image
            className='mx-auto w-auto'
            alt='Logo'
            width={180}
            height={180}
            {...props}
        />
    );
}
