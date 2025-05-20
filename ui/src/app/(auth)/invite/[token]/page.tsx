'use client';
import {useState} from 'react';
import {useRouter} from 'next/navigation';

interface InviteTokenPageProps {
    params: {token: string};
}

export default function InviteTokenPage({params}: InviteTokenPageProps) {
    const router = useRouter();
    const token = params.token;

    const [name, setName] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        setLoading(true);
        try {
            const res = await fetch('/api/invite/accept', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({token, name, password}),
            });
            if (!res.ok) {
                const data = await res.json();
                setError(data.error || 'Registration failed');
                setLoading(false);
                return;
            }
            router.replace('/login?registered=1');
        } catch {
            setError('Network error');
            setLoading(false);
        }
    };

    return (
        <div className='max-w-md mx-auto mt-16 p-8 bg-white rounded shadow'>
            <h1 className='text-2xl font-bold mb-4'>Accept Invitation</h1>
            <form onSubmit={handleSubmit} className='space-y-4'>
                <div>
                    <label className='block mb-1 font-medium'>Name</label>
                    <input
                        type='text'
                        className='w-full border rounded px-3 py-2'
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        required
                    />
                </div>
                <div>
                    <label className='block mb-1 font-medium'>Password</label>
                    <input
                        type='password'
                        className='w-full border rounded px-3 py-2'
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        required
                    />
                </div>
                {error && <div className='text-red-600'>{error}</div>}
                <button
                    type='submit'
                    className='w-full bg-blue-600 text-white py-2 rounded font-semibold hover:bg-blue-700'
                    disabled={loading}
                >
                    {loading ? 'Registering...' : 'Register'}
                </button>
            </form>
        </div>
    );
}
