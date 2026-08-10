import { Toaster } from "@/components/ui/sonner"
import { Providers } from "@/components/providers";
import type { Metadata } from "next";
import { getBaseUrl } from "@/lib/base-url";
import { cn } from "@/lib/utils";
import "./globals.css";

const appUrl = getBaseUrl();
const socialImage = "/images/social-card.png";
const description = "The No-Hassle CMS for GitHub";

export const metadata: Metadata = {
  metadataBase: new URL(appUrl),
  title: {
    template: "%s | Pages CMS",
    default: "Pages CMS",
  },
  description,
  alternates: {
    canonical: "/",
  },
  openGraph: {
    type: "website",
    url: appUrl,
    siteName: "Pages CMS",
    title: "Pages CMS",
    description,
    images: [
      {
        url: socialImage,
        width: 1200,
        height: 630,
        alt: "Pages CMS social card",
      },
    ],
  },
  twitter: {
    card: "summary_large_image",
    title: "Pages CMS",
    description,
    images: [socialImage],
  },
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
	return (
    <html lang="en" suppressHydrationWarning>
      <body
        className={cn(
          "min-h-screen bg-background font-sans antialiased",
        )}
      >
        <Providers user={null}>
          {children}
        </Providers>
        <Toaster/>
      </body>
    </html>
  );
}
