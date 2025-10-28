import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'
import { DarkModeProvider } from '@/contexts/DarkModeContext'
import Navbar from '@/components/Navbar'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'EstateSaleFinder.ai - Find Estate Sales Near You',
  description: 'Discover estate sales, auctions, and moving sales in your area. AI-powered search and personalized recommendations.',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <DarkModeProvider>
          <Navbar />
          {children}
        </DarkModeProvider>
      </body>
    </html>
  )
}
