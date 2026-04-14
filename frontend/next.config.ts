import type { NextConfig } from "next";
import { initOpenNextCloudflareForDev } from "@opennextjs/cloudflare";

const nextConfig: NextConfig = {
	turbopack: {
		resolveAlias: {
			"@/*": "./src/*",
		},
	},
};

initOpenNextCloudflareForDev();

export default nextConfig;
