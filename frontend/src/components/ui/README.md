# VamsaSetu UI Components

A collection of reusable, accessible UI components built with React, TypeScript, Tailwind CSS, and Framer Motion.

## Components

### Button

A versatile button component with multiple variants, sizes, and states.

**Props:**
- `variant`: 'primary' | 'secondary' | 'outline' (default: 'primary')
- `size`: 'sm' | 'md' | 'lg' (default: 'md')
- `isLoading`: boolean (default: false)
- `leftIcon`: ReactNode
- `rightIcon`: ReactNode
- `fullWidth`: boolean (default: false)
- All standard HTML button attributes

**Usage:**
```tsx
import { Button } from '@/components/ui';

<Button variant="primary" onClick={handleClick}>
  Click Me
</Button>

<Button variant="secondary" leftIcon={<Icon />} isLoading>
  Loading...
</Button>
```

**Variants:**
- **Primary**: Saffron background with white text
- **Secondary**: Teal background with white text
- **Outline**: Transparent with saffron border

### Input

A form input component with label, error display, and validation support.

**Props:**
- `label`: string
- `error`: string
- `helperText`: string
- `leftIcon`: ReactNode
- `rightIcon`: ReactNode
- `fullWidth`: boolean (default: false)
- All standard HTML input attributes

**Usage:**
```tsx
import { Input } from '@/components/ui';

<Input
  label="Email"
  type="email"
  placeholder="your@email.com"
  error={errors.email}
  required
  fullWidth
/>

<Input
  label="Search"
  leftIcon={<SearchIcon />}
  placeholder="Search..."
/>
```

**Features:**
- Automatic error styling
- Required field indicator
- Helper text support
- Icon support (left/right)
- Accessible with ARIA attributes

### Modal

A modal dialog component with Framer Motion animations and portal rendering.

**Props:**
- `isOpen`: boolean (required)
- `onClose`: () => void (required)
- `title`: string
- `children`: ReactNode (required)
- `size`: 'sm' | 'md' | 'lg' | 'xl' (default: 'md')
- `showCloseButton`: boolean (default: true)
- `closeOnOverlayClick`: boolean (default: true)
- `closeOnEscape`: boolean (default: true)

**Usage:**
```tsx
import { Modal } from '@/components/ui';

const [isOpen, setIsOpen] = useState(false);

<Modal
  isOpen={isOpen}
  onClose={() => setIsOpen(false)}
  title="Add Member"
  size="md"
>
  <div>Modal content here</div>
</Modal>
```

**Features:**
- Smooth enter/exit animations
- Backdrop blur effect
- Keyboard navigation (Escape to close)
- Body scroll lock when open
- Portal rendering (renders at document.body)
- Accessible with ARIA attributes

### Card

A flexible card component with sub-components for structured content.

**Props:**
- `children`: ReactNode (required)
- `className`: string
- `variant`: 'default' | 'elevated' | 'outlined' (default: 'default')
- `padding`: 'none' | 'sm' | 'md' | 'lg' (default: 'md')
- `hoverable`: boolean (default: false)
- `onClick`: () => void

**Sub-components:**
- `CardHeader`: Container for card header content
- `CardTitle`: Styled title text
- `CardDescription`: Styled description text
- `CardContent`: Main content area
- `CardFooter`: Footer with top border

**Usage:**
```tsx
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui';

<Card variant="elevated">
  <CardHeader>
    <CardTitle>Card Title</CardTitle>
    <CardDescription>Card description</CardDescription>
  </CardHeader>
  <CardContent>
    <p>Main content here</p>
  </CardContent>
  <CardFooter>
    <Button>Action</Button>
  </CardFooter>
</Card>

<Card hoverable onClick={handleClick}>
  <CardTitle>Clickable Card</CardTitle>
</Card>
```

**Variants:**
- **Default**: Light shadow with border
- **Elevated**: Larger shadow, no border
- **Outlined**: Saffron border, no shadow

## Theme Colors

All components use VamsaSetu theme colors defined in `src/index.css`:

- **Saffron** (#E8650A): Primary brand color
- **Saffron Light** (#F5A623): Hover states
- **Ivory** (#FBF5E6): Background
- **Teal** (#0D4A52): Secondary actions
- **Teal Light** (#0D9488): Teal hover states
- **Charcoal** (#2C2420): Text
- **Rose** (#E11D48): Errors and destructive actions
- **Amber** (#F59E0B): Accents

## Typography

- **Headings**: Playfair Display (serif)
- **Body**: DM Sans (sans-serif)

## Accessibility

All components follow accessibility best practices:

- Semantic HTML elements
- ARIA attributes where appropriate
- Keyboard navigation support
- Focus indicators
- Screen reader support
- Minimum touch target size (44px)

## Animation

Components use Framer Motion for smooth animations:

- Button: Scale on hover/tap
- Modal: Fade and scale entrance
- Card: Lift on hover (when hoverable)

## Examples

See `UIComponents.example.tsx` for complete usage examples including:
- Button variants and states
- Input validation patterns
- Modal dialogs
- Card layouts
- Complete form example

## Testing

Unit tests are provided for all components:
- `Button.test.tsx`
- `Input.test.tsx`
- `Modal.test.tsx`
- `Card.test.tsx`

Run tests with: `npm test` (when test setup is configured)
