import { createMDX } from "fumadocs-mdx/next";

const basePath = process.env.PAGES_BASE_PATH ?? "";

const config = {
  reactStrictMode: true,
  output: "export",
  basePath,
  assetPrefix: basePath || undefined,
  images: { unoptimized: true },
  trailingSlash: true,
};

export default createMDX()(config);
