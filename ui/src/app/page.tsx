import Link from 'next/link';

export default function HomePage() {
    return (
        <div className='max-w-md mx-auto mt-12 p-8 bg-white rounded-xl shadow-md border border-gray-200 text-center'>
            <h1 className='text-3xl font-bold mb-6'>WebAuthn Demo</h1>
            <div className='mt-8 flex flex-col gap-4'>
                <Link
                    href='/login'
                    className='bg-green-500 hover:bg-green-600 text-white font-semibold py-2 px-4 rounded transition-colors'
                >
                    Login
                </Link>
                <Link
                    href='/register'
                    className='bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded transition-colors'
                >
                    Register
                </Link>
            </div>
        </div>
    );
}
