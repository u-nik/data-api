'use client';

import {useState} from 'react';
import {webauthnRegister} from './WebAuthn';
import {Heading} from './catalyst/heading';
import {Field, Label} from './catalyst/fieldset';
import {Input} from './catalyst/input';
import {Button} from './catalyst/button';

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
        <div className='grid w-full max-w-sm grid-cols-1 gap-8'>
            <Heading>Register with Passkey</Heading>
            <Field>
                <Label>Username</Label>
                <Input
                    placeholder='Username'
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    autoComplete='username webauthn'
                    autoFocus={true}
                ></Input>
            </Field>
            <Button type='submit' className='w-full' onClick={handleRegister}>
                Register
            </Button>

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
