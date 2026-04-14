/** @type {import('next').NextConfig} */
const nextConfig = {}

export default nextConfig

import('@opennextjs/cloudflare').then(m => m.initOpenNextCloudflareForDev());
