import { AuthProvider } from '@/providers/auth-context'
import '@/globals.css'

export const metadata = {
  title: 'Owndangan - Wedding Invitation Platform',
  description: 'Create beautiful wedding invitations online',
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>
        <AuthProvider>
          {children}
        </AuthProvider>
      </body>
    </html>
  )
}
