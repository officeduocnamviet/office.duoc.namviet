import { withSentryConfig } from "@sentry/nextjs";
import type { NextConfig } from "next";

const withPWA = require('next-pwa')({
  dest: 'public',
  disable: process.env.NODE_ENV === 'development',
});

const nextConfig: NextConfig = {
  /* config options here */
};

export default withSentryConfig(
  withPWA(nextConfig),
  {
    silent: true,
    org: "nam-viet-erp",
    project: "web-admin",
    widenClientFileUpload: true,
    tunnelRoute: "/monitoring",
  }
);
