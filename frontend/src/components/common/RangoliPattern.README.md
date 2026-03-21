# RangoliPattern Component

A decorative SVG background component inspired by traditional Indian rangoli designs, featuring geometric patterns with symmetry using VamsaSetu theme colors.

**Validates: Requirements 10.3**

## Overview

The RangoliPattern component creates a beautiful, culturally-inspired background texture that can be used throughout the VamsaSetu application. It features:

- **Traditional Design**: Lotus-inspired mandala pattern with radial symmetry
- **Theme Colors**: Uses VamsaSetu's saffron (#E8650A), turmeric (#F5A623), and teal (#0D4A52)
- **SVG-based**: Crisp rendering at any scale
- **Customizable**: Adjustable opacity and pattern size
- **Performance**: Lightweight and non-interactive (pointer-events-none)
- **Accessible**: Properly marked with aria-hidden

## Design Elements

The rangoli pattern includes:

1. **Central Mandala**: 8-petaled lotus design with layered petals
2. **Color Layers**:
   - Outer petals: Saffron (primary brand color)
   - Middle petals: Turmeric (warm accent)
   - Inner circle: Teal (depth and contrast)
   - Center dot: Saffron (focal point)
3. **Corner Decorations**: Small circular elements at pattern tile corners
4. **Geometric Grid**: Subtle connecting lines for structure

## Usage

### Basic Usage

```tsx
import { RangoliPattern } from '@/components/common/RangoliPattern';

function MyPage() {
  return (
    <div className="relative min-h-screen">
      <RangoliPattern />
      <div className="relative z-10">
        {/* Your content here */}
      </div>
    </div>
  );
}
```

### With Custom Opacity

```tsx
// More visible pattern
<RangoliPattern opacity={0.2} />

// Very subtle pattern
<RangoliPattern opacity={0.05} />
```

### With Custom Pattern Size

```tsx
// Larger tiles for spacious feel
<RangoliPattern size={300} />

// Smaller tiles for denser texture
<RangoliPattern size={150} />
```

### With Custom Styling

```tsx
<RangoliPattern 
  opacity={0.15} 
  className="mix-blend-multiply" 
/>
```

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `opacity` | `number` | `0.1` | Opacity of the pattern (0-1). Lower values are more subtle. |
| `size` | `number` | `200` | Size of the pattern tile in pixels. Affects pattern density. |
| `className` | `string` | `''` | Additional CSS classes for custom styling. |

## Recommended Use Cases

### 1. Page Backgrounds
Use with low opacity (0.08-0.12) for full-page backgrounds:

```tsx
<div className="relative min-h-screen bg-ivory">
  <RangoliPattern opacity={0.1} size={250} />
  <div className="relative z-10">
    {/* Page content */}
  </div>
</div>
```

### 2. Card Backgrounds
Add subtle texture to cards and panels:

```tsx
<div className="relative bg-white rounded-lg p-6 overflow-hidden">
  <RangoliPattern opacity={0.12} />
  <div className="relative z-10">
    {/* Card content */}
  </div>
</div>
```

### 3. Hero Sections
Create visual interest in landing sections:

```tsx
<section className="relative py-20 bg-gradient-to-br from-ivory to-white">
  <RangoliPattern opacity={0.15} size={300} />
  <div className="relative z-10 container mx-auto">
    {/* Hero content */}
  </div>
</section>
```

### 4. Modal Backgrounds
Add depth to modal dialogs:

```tsx
<div className="relative bg-white rounded-xl p-8 overflow-hidden">
  <RangoliPattern opacity={0.08} size={180} />
  <div className="relative z-10">
    {/* Modal content */}
  </div>
</div>
```

## Design Guidelines

### Opacity Recommendations

- **0.05-0.08**: Very subtle, for text-heavy pages
- **0.1-0.12**: Default subtle, for general backgrounds
- **0.15-0.2**: More visible, for decorative sections
- **0.2+**: Prominent, use sparingly for special emphasis

### Size Recommendations

- **150px**: Dense texture, compact sections
- **200px**: Default balanced, general use
- **250-300px**: Spacious, hero sections
- **350px+**: Very spacious, large displays

### Best Practices

1. **Always use relative positioning** on the parent container
2. **Add `overflow-hidden`** to prevent pattern overflow
3. **Use `z-10` or higher** on content to ensure it appears above the pattern
4. **Keep opacity low** (≤0.15) for readability
5. **Test on different backgrounds** - works best on white/ivory
6. **Consider performance** - one pattern per viewport is usually sufficient

## Accessibility

The component is properly configured for accessibility:

- `aria-hidden="true"`: Hidden from screen readers (decorative only)
- `pointer-events-none`: Doesn't interfere with user interactions
- Non-interactive: Purely visual decoration

## Cultural Context

Rangoli is a traditional Indian art form where patterns are created on the ground using colored powders, flowers, or rice. These designs are:

- **Auspicious**: Used during festivals and celebrations
- **Welcoming**: Placed at entrances to welcome guests
- **Geometric**: Feature symmetry and mathematical precision
- **Colorful**: Use vibrant, meaningful colors

The VamsaSetu rangoli pattern honors this tradition while maintaining a modern, digital aesthetic suitable for a family tree application.

## Technical Details

- **SVG-based**: Vector graphics scale perfectly at any resolution
- **Unique IDs**: Each instance generates a unique pattern ID to avoid conflicts
- **Absolute positioning**: Positioned absolutely within relative parent
- **No JavaScript**: Pure CSS/SVG, no runtime overhead
- **Lightweight**: Minimal DOM nodes, efficient rendering

## Examples

See `RangoliPattern.example.tsx` for comprehensive usage examples including:
- Default subtle backgrounds
- Various opacity levels
- Different pattern sizes
- Full-page layouts
- Custom blend modes
- Integration with gradients

## Related Components

- `Layout`: Main layout component that could use RangoliPattern
- `Card`: UI component that pairs well with subtle patterns
- `Modal`: Dialog component that can use pattern backgrounds

## Version History

- **v1.0.0**: Initial implementation with lotus mandala design
  - 8-petaled symmetrical pattern
  - VamsaSetu theme colors
  - Customizable opacity and size
  - Accessibility features
