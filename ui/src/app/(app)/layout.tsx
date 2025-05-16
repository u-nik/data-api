import {Link} from '@/components/ui/link';
import {Logo} from '@/components/ui/logo';
import {Navbar, NavbarItem, NavbarSection} from '@/components/ui/navbar';
import {StackedLayout} from '@/components/ui/stacked-layout';

export default function RootLayout({
    children,
}: Readonly<{children: React.ReactNode}>) {
    return (
        <StackedLayout
            navbar={
                <Navbar>
                    <Link href='/' aria-label='Home'>
                        <Logo src='/logo.png' className='size-10 sm:size-8' />
                    </Link>
                    <NavbarSection>
                        <NavbarItem href='/' current>
                            Home
                        </NavbarItem>
                        <NavbarItem href='/events'>Events</NavbarItem>
                        <NavbarItem href='/orders'>Orders</NavbarItem>
                    </NavbarSection>
                </Navbar>
            }
            sidebar={null}
        >
            {children}
        </StackedLayout>
    );
}
