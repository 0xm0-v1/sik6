# Web Image Standards

This guide defines how to use images in our web project.  
It covers formats, sizes, and usage rules to ensure **performance, quality, and compatibility**.

## Image Formats

| Format   | When to Use                                               | Notes                                                                       |
| -------- | --------------------------------------------------------- | --------------------------------------------------------------------------- |
| **SVG**  | Logos, icons, simple graphics                             | Vector, scalable, lightweight. Preferred whenever possible.                 |
| **WebP** | General-purpose images (UI, content, backgrounds)         | Best compromise quality/size. Supported by all modern browsers.             |
| **AVIF** | Large media, photos, hero images                          | Smaller than WebP, excellent quality, but slower encoding. Fallback needed. |
| **PNG**  | When transparency is required **and** SVG is not possible | Lossless, heavier. Use only if WebP/AVIF not viable.                        |
| **JPG**  | Photos when legacy support is required                    | Use progressive compression. Deprecated in most new use cases.              |

> **Default format: WebP**, with fallback to PNG/JPG if legacy support is critical.

## Responsive Images

Always provide multiple resolutions and let the browser pick:

```html
<img
  src="image-800.webp"
  srcset="image-400.webp 400w, image-800.webp 800w, image-1600.webp 1600w"
  sizes="(max-width: 600px) 100vw, 600px"
  alt="Description"
/>
```

- Provide `1×`, `2×`, `3×` versions for Retina/HiDPI screens.
- Use `loading="lazy"` for below-the-fold images.
- Compress aggressively (but avoid visible artifacts).

## Special Cases & Rules of Thumb

| Category          | Specification / Rule                      | Notes                             |
| ----------------- | ----------------------------------------- | --------------------------------- |
| **Favicons**      | Multi-size `.ico` or PNGs                 | Sizes: 16, 32, 48, 64 px          |
| **PWA / App**     | Icons: 512×512, 192×192, 180×180, 256×256 | Use PNG or WebP                   |
| **Open Graph**    | Facebook/LinkedIn → 1200×630 px (1.91:1)  | Must be rectangular               |
| **Open Graph**    | Twitter/X → 1200×600 px (2:1)             | Do not stretch square logos       |
| **Logos**         | Use **SVG** whenever possible             | Scalable, minimal                 |
| **Bitmap**        | Default: **WebP**                         | Best balance quality/performance  |
| **Fallbacks**     | PNG/JPG only if WebP/AVIF not supported   | Avoid as primary format           |
| **Optimization**  | Always compress images                    | Tools: ImageOptim, Squoosh, Sharp |
| **Assets**        | Do **not** upload raw design exports      | Avoid uncompressed PNGs           |
| **Hero images**   | Prefer WebP/AVIF with `srcset`            | Enable responsive loading         |
| **Accessibility** | Always include **meaningful alt text**    | SEO + screen readers              |

_This system ensures that images are lightweight, responsive, and accessible, reducing load times and improving the overall user experience while maintaining compatibility across modern browsers._
