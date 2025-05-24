import {NextRequest, NextResponse} from 'next/server';

export async function POST(req: NextRequest) {
    const {token, name, password} = await req.json();
    if (!token || !name || !password) {
        return NextResponse.json({error: 'Missing fields'}, {status: 400});
    }

    // Call backend API to accept invitation
    const res = await fetch(
        process.env.DATA_API_URL + '/api/invitations/accept',
        {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({token, name, password}),
        },
    );

    if (!res.ok) {
        const data = await res.json().catch(() => ({}));
        return NextResponse.json(
            {
                error: data.error || 'Failed to accept invitation',
                result: data.result || null,
            },
            {status: 400},
        );
    }

    return NextResponse.json({success: true});
}
