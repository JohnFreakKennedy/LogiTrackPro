# LogiTrackPro Design Specification

## Document Overview

This document provides comprehensive design specifications for the LogiTrackPro logistics planning platform. It covers visual design, user interface components, interaction patterns, and design system guidelines.

**Version:** 1.0  
**Last Updated:** 2024  
**Design System:** Custom Dark Theme with Green/Orange Accents

---

## 1. Design Philosophy

### Core Principles
- **Clarity First**: Information hierarchy prioritizes actionable logistics data
- **Efficiency**: Minimize clicks and navigation for common tasks
- **Visual Feedback**: Clear status indicators for optimization processes
- **Responsive**: Seamless experience across desktop, tablet, and mobile
- **Accessibility**: WCAG 2.1 AA compliance for color contrast and interactions

### Visual Identity
- **Primary Color**: Green (#22c55e) - Represents optimization, success, efficiency
- **Accent Color**: Orange (#f97316) - Highlights important actions and alerts
- **Dark Theme**: Modern dark interface reducing eye strain for extended use
- **Typography**: DM Sans (body), Outfit (headings), JetBrains Mono (data)

---

## 2. Color System

### Primary Palette
```
Primary Green:
- 50:  #f0fdf4  (Lightest)
- 100: #dcfce7
- 200: #bbf7d0
- 300: #86efac
- 400: #4ade80
- 500: #22c55e  (Base - Primary Actions)
- 600: #16a34a  (Hover States)
- 700: #15803d
- 800: #166534
- 900: #14532d  (Darkest)
```

### Accent Palette
```
Orange Accent:
- Light: #fb923c
- Base:  #f97316  (Primary Accent)
- Dark:  #ea580c
```

### Dark Theme Palette
```
Background:
- 50:   #f8fafc  (Text on Dark)
- 100:  #f1f5f9
- 200:  #e2e8f0
- 300:  #cbd5e1
- 400:  #94a3b8  (Secondary Text)
- 500:  #64748b
- 600:  #475569
- 700:  #334155  (Borders)
- 800:  #1e293b  (Cards, Surfaces)
- 900:  #0f172a  (Main Background)
- 950:  #020617  (Deep Background)
```

### Semantic Colors
```
Success:  #22c55e (Primary-500)
Warning:  #f59e0b (Yellow-500)
Error:    #ef4444 (Red-500)
Info:     #3b82f6 (Blue-500)
```

---

## 3. Typography

### Font Families
- **Body Text**: DM Sans (400, 500, 600, 700)
- **Headings**: Outfit (400, 500, 600, 700, 800)
- **Data/Monospace**: JetBrains Mono (400, 500)

### Type Scale
```
Display:  3rem (48px)  - Hero Headings
H1:       2.25rem (36px) - Page Titles
H2:       1.875rem (30px) - Section Headers
H3:       1.5rem (24px) - Subsection Headers
H4:       1.25rem (20px) - Card Titles
Body:     1rem (16px) - Default Text
Small:    0.875rem (14px) - Secondary Text
Tiny:     0.75rem (12px) - Labels, Captions
```

### Font Weights
- **Regular**: 400 (Body text)
- **Medium**: 500 (Emphasized text, buttons)
- **Semibold**: 600 (Headings, labels)
- **Bold**: 700 (Strong emphasis)
- **Extra Bold**: 800 (Display text)

---

## 4. Layout System

### Grid System
- **Container Max Width**: 1440px
- **Grid Columns**: 12-column responsive grid
- **Gutter**: 24px (desktop), 16px (tablet), 12px (mobile)
- **Breakpoints**:
  - Mobile: < 640px
  - Tablet: 640px - 1024px
  - Desktop: > 1024px

### Spacing Scale
```
xs:   4px   (0.25rem)
sm:   8px   (0.5rem)
md:   16px  (1rem)
lg:   24px  (1.5rem)
xl:   32px  (2rem)
2xl:  48px  (3rem)
3xl:  64px  (4rem)
```

### Layout Components

#### Sidebar Navigation
- **Width**: 256px (16rem) fixed on desktop
- **Background**: dark-900/80 with backdrop blur
- **Border**: Right border dark-700/50
- **Height**: Full viewport height
- **Mobile**: Collapsible overlay menu

#### Main Content Area
- **Padding**: 32px (2rem) desktop, 16px (1rem) mobile
- **Margin Left**: 256px on desktop (sidebar width)
- **Background**: Gradient from dark-900 to dark-950

---

## 5. Component Library

### 5.1 Buttons

#### Primary Button
```
Background: primary-600
Hover: primary-500
Text: white
Shadow: shadow-lg shadow-primary-600/25
Padding: 16px 24px (px-4 py-2)
Border Radius: 8px (rounded-lg)
Font: Medium (500)
```

#### Secondary Button
```
Background: dark-700
Hover: dark-600
Text: white
Border: 1px solid dark-600
Padding: 16px 24px
Border Radius: 8px
```

#### Accent Button
```
Background: accent (#f97316)
Hover: accent-light (#fb923c)
Text: white
Shadow: shadow-lg shadow-accent/25
```

#### Danger Button
```
Background: red-600
Hover: red-500
Text: white
```

### 5.2 Cards

#### Standard Card
```
Background: dark-800/50 with backdrop-blur-sm
Border: 1px solid dark-700/50
Border Radius: 12px (rounded-xl)
Padding: 24px (p-6)
Shadow: Subtle inner glow
```

#### Interactive Card (Hover)
```
Hover Border: primary-500/30
Hover Shadow: shadow-lg shadow-primary-500/5
Transition: all 300ms ease
```

### 5.3 Forms

#### Input Fields
```
Background: dark-800 (dark-700 on focus)
Border: 1px solid dark-700
Focus Border: primary-500
Focus Shadow: 0 0 0 3px rgba(34, 197, 94, 0.1)
Border Radius: 8px
Padding: 12px 16px (0.75rem 1rem)
Text Color: white
Placeholder: dark-400
```

#### Select Dropdowns
- Same styling as inputs
- Custom dropdown arrow (chevron-down icon)
- Options background: dark-800
- Hover state: dark-700

#### Textarea
- Same base styling as inputs
- Min height: 100px
- Resizable: vertical only

### 5.4 Tables

#### Table Container
```
Background: dark-800/50
Border: 1px solid dark-700/50
Border Radius: 12px
Overflow: Auto (horizontal scroll on mobile)
```

#### Table Header (th)
```
Background: dark-800
Text Color: dark-300
Font: Medium (500), 14px
Text Transform: Uppercase
Letter Spacing: 0.05em
Padding: 12px 16px
```

#### Table Cell (td)
```
Padding: 12px 16px
Border Top: 1px solid dark-700/50
Text Color: white
```

#### Table Row Hover
```
Background: dark-800/30
Transition: background 200ms ease
```

### 5.5 Badges

#### Success Badge
```
Background: primary-500/20
Text: primary-400
Border Radius: 9999px (rounded-full)
Padding: 4px 10px
Font: 12px, Medium
```

#### Warning Badge
```
Background: yellow-500/20
Text: yellow-400
```

#### Info Badge
```
Background: blue-500/20
Text: blue-400
```

#### Danger Badge
```
Background: red-500/20
Text: red-400
```

### 5.6 Modals

#### Modal Backdrop
```
Background: black/60 with backdrop-blur-sm
Position: Fixed, full viewport
Z-index: 50
```

#### Modal Container
```
Background: dark-800
Border: 1px solid dark-700/50
Border Radius: 16px (rounded-2xl)
Shadow: shadow-2xl
Max Width: 512px (max-w-lg)
Max Height: 90vh
Overflow: Hidden
```

#### Modal Header
```
Padding: 24px (p-6)
Border Bottom: 1px solid dark-700/50
Display: Flex, justify-between, items-center
```

#### Modal Content
```
Padding: 24px (p-6)
Max Height: calc(90vh - 80px)
Overflow: Auto
```

### 5.7 Navigation

#### Sidebar Navigation Item
```
Inactive:
- Background: Transparent
- Text: dark-300
- Hover Background: dark-800
- Hover Text: white

Active:
- Background: primary-500/10
- Text: primary-400
- Border: 1px solid primary-500/20
```

#### Mobile Menu
```
Background: dark-900/95 with backdrop-blur-xl
Position: Fixed, full width
Border Bottom: 1px solid dark-700/50
Z-index: 50
```

---

## 6. Page Specifications

### 6.1 Login Page

#### Layout
- **Centered Card**: Max width 400px, centered vertically and horizontally
- **Background**: Dark gradient (dark-900 to dark-950)
- **Card**: Standard card styling with extra padding

#### Components
- Logo: 48px × 48px gradient icon (primary-500 to primary-600)
- Title: "Welcome to LogiTrackPro" (H1, Outfit Bold)
- Subtitle: "Sign in to your account" (Body, dark-400)
- Email Input: Full width
- Password Input: Full width with show/hide toggle
- Submit Button: Primary button, full width
- Register Link: Centered below form, accent color

#### States
- **Loading**: Button shows spinner, disabled state
- **Error**: Red error message below form
- **Success**: Redirect to dashboard

### 6.2 Dashboard

#### Header Section
- **Title**: "Dashboard" (H1, Outfit Bold)
- **Subtitle**: "Overview of your logistics operations" (Body, dark-400)

#### Stats Grid (4 columns)
Each stat card:
- **Icon**: 48px × 48px gradient circle (color varies)
- **Label**: dark-400, 14px
- **Value**: 36px, Outfit Bold, white
- **Link**: "View all" with arrow icon, hover effect
- **Colors**: Primary (warehouses), Blue (customers), Purple (vehicles), Accent (plans)

#### Performance Metrics Card
- **Title**: "Performance Metrics" with TrendingUp icon
- **Metrics**: 3 items in vertical list
  - Icon + Label + Value (monospace font)
  - Background: dark-800/50
  - Padding: 16px
  - Border Radius: 12px

#### Recent Plans Card (2 columns wide)
- **Header**: Title + "View all" link
- **Plan Items**: List of plan cards
  - Plan name (bold, hover: primary-400)
  - Date range (small, dark-400)
  - Cost (monospace, right-aligned)
  - Status badge
- **Empty State**: Centered icon + message + CTA link

#### Quick Actions Card
- **Title**: "Quick Actions"
- **Grid**: 4 buttons (2 columns mobile, 4 columns desktop)
- **Buttons**: Secondary style with icons

### 6.3 Master Data Pages (Warehouses, Customers, Vehicles)

#### Common Layout
- **Header**: Title + "Add New" button (primary)
- **Table**: Full-width table with columns:
  - Name/Identifier
  - Location/Details
  - Capacity/Attributes
  - Status/Availability
  - Actions (Edit, Delete)

#### Add/Edit Modal
- **Form Fields**: Based on entity type
- **Layout**: Single column, stacked
- **Actions**: Cancel (secondary) + Save (primary), right-aligned

#### Warehouse Specific
- Fields: Name, Address, Coordinates (lat/lng), Capacity, Current Stock, Holding Cost, Replenishment Qty

#### Customer Specific
- Fields: Name, Address, Coordinates, Demand Rate, Max/Min Inventory, Current Inventory, Holding Cost, Priority

#### Vehicle Specific
- Fields: Name, Capacity, Cost per km, Fixed Cost, Max Distance, Warehouse (dropdown), Available (toggle)

### 6.4 Plans Page

#### Header
- **Title**: "Delivery Plans"
- **Actions**: "Create Plan" button (accent)

#### Plans List
- **Card Layout**: Grid (1 column mobile, 2-3 columns desktop)
- **Card Content**:
  - Plan name (bold)
  - Date range
  - Status badge
  - Metrics: Cost, Distance, Routes count
  - Actions: View Details, Optimize, Delete

#### Create Plan Modal
- **Fields**: Name, Start Date, End Date, Warehouse (dropdown)
- **Validation**: Date range validation, warehouse required

### 6.5 Plan Detail Page

#### Header Section
- **Plan Name**: Large heading
- **Status Badge**: Prominent display
- **Actions**: Optimize button (accent), Delete button (danger)

#### Plan Information Card
- **Fields**: Name, Dates, Warehouse, Status, Total Cost, Total Distance
- **Layout**: Grid, 2 columns

#### Routes Section
- **Day Tabs**: Horizontal tabs for each day
- **Route Cards**: Per vehicle
  - Vehicle name
  - Distance, Cost, Load
  - Stops list with sequence
  - Map visualization (future)

#### Optimization Status
- **Loading State**: Spinner + "Optimizing routes..." message
- **Success State**: Green checkmark + summary
- **Error State**: Red alert + error message

---

## 7. Interaction Patterns

### 7.1 Animations

#### Page Transitions
```
Initial: opacity: 0, y: 20px
Animate: opacity: 1, y: 0
Duration: 300ms
Easing: ease-out
```

#### Card Hover
```
Scale: 1.02 (subtle)
Shadow: Increase intensity
Border: Highlight color
Duration: 200ms
```

#### Button Press
```
Scale: 0.98
Duration: 100ms
```

#### Modal Entrance
```
Initial: opacity: 0, scale: 0.95, y: 20px
Animate: opacity: 1, scale: 1, y: 0
Duration: 200ms
Easing: ease-out
```

### 7.2 Loading States

#### Spinner
- **Size**: 32px × 32px
- **Color**: primary-500
- **Animation**: Rotate 360deg, infinite, 1s linear
- **Usage**: Full page, buttons, inline

#### Skeleton Loaders
- **Background**: dark-800
- **Shimmer**: Animated gradient overlay
- **Usage**: Table rows, cards, lists

### 7.3 Feedback States

#### Success
- **Color**: primary-500
- **Icon**: Checkmark circle
- **Message**: Green text, clear action confirmation

#### Error
- **Color**: red-500
- **Icon**: Alert circle
- **Message**: Red text, actionable error message

#### Warning
- **Color**: yellow-500
- **Icon**: Alert triangle
- **Message**: Yellow text, cautionary information

#### Info
- **Color**: blue-500
- **Icon**: Info circle
- **Message**: Blue text, informational content

---

## 8. Responsive Design

### 8.1 Mobile (< 640px)
- **Sidebar**: Hidden, replaced with hamburger menu
- **Grid**: Single column layouts
- **Tables**: Horizontal scroll or card layout
- **Padding**: Reduced (12px-16px)
- **Typography**: Slightly smaller scale
- **Buttons**: Full width in forms

### 8.2 Tablet (640px - 1024px)
- **Sidebar**: Collapsible or overlay
- **Grid**: 2 columns for stats
- **Tables**: Horizontal scroll enabled
- **Modals**: Max width 90vw

### 8.3 Desktop (> 1024px)
- **Sidebar**: Fixed, always visible
- **Grid**: Full 4-column layouts
- **Tables**: Full width, no scroll
- **Modals**: Max width 512px

---

## 9. Accessibility

### 9.1 Color Contrast
- **Text on Dark**: Minimum 4.5:1 contrast ratio
- **Interactive Elements**: Minimum 3:1 contrast ratio
- **Focus Indicators**: 2px solid primary-500 outline

### 9.2 Keyboard Navigation
- **Tab Order**: Logical flow through interactive elements
- **Focus Visible**: Clear focus indicators on all focusable elements
- **Skip Links**: Skip to main content link

### 9.3 Screen Readers
- **ARIA Labels**: All icons and interactive elements
- **Semantic HTML**: Proper heading hierarchy
- **Alt Text**: Descriptive text for all images/icons

### 9.4 Motion
- **Respect Preferences**: Honor `prefers-reduced-motion`
- **Animation Duration**: Maximum 300ms for transitions

---

## 10. Iconography

### Icon Library
- **Source**: Lucide React Icons
- **Size**: 16px (small), 20px (default), 24px (large), 32px (xlarge)
- **Color**: Inherit from parent or explicit color classes
- **Stroke Width**: 2px default

### Common Icons
- **Dashboard**: LayoutDashboard
- **Warehouses**: Warehouse
- **Customers**: Users
- **Vehicles**: Truck
- **Plans**: Route
- **Settings**: Settings
- **Logout**: LogOut
- **Add**: Plus
- **Edit**: Pencil
- **Delete**: Trash
- **Search**: Search
- **Filter**: Filter
- **Export**: Download
- **Loading**: Loader2 (spinning)

---

## 11. Data Visualization (Future)

### 11.1 Charts
- **Library**: Recharts or Chart.js
- **Colors**: Primary palette for data series
- **Background**: dark-800
- **Grid Lines**: dark-700
- **Text**: white/dark-300

### 11.2 Maps
- **Provider**: Mapbox or Leaflet
- **Style**: Dark theme map
- **Markers**: Custom SVG icons
- **Routes**: Primary-500 colored lines
- **Warehouses**: Green markers
- **Customers**: Blue markers

---

## 12. Design Tokens

### Spacing Tokens
```javascript
spacing: {
  xs: '4px',
  sm: '8px',
  md: '16px',
  lg: '24px',
  xl: '32px',
  '2xl': '48px',
  '3xl': '64px'
}
```

### Border Radius Tokens
```javascript
radius: {
  sm: '4px',
  md: '8px',
  lg: '12px',
  xl: '16px',
  full: '9999px'
}
```

### Shadow Tokens
```javascript
shadows: {
  sm: '0 1px 2px rgba(0, 0, 0, 0.05)',
  md: '0 4px 6px rgba(0, 0, 0, 0.1)',
  lg: '0 10px 15px rgba(0, 0, 0, 0.1)',
  xl: '0 20px 25px rgba(0, 0, 0, 0.15)',
  '2xl': '0 25px 50px rgba(0, 0, 0, 0.25)',
  glow: '0 0 40px rgba(34, 197, 94, 0.15)'
}
```

---

## 13. Implementation Notes

### 13.1 CSS Framework
- **Base**: Tailwind CSS 3.4+
- **Custom Config**: Extended theme with design tokens
- **PostCSS**: Autoprefixer for browser compatibility

### 13.2 Animation Library
- **Library**: Framer Motion
- **Usage**: Page transitions, modal animations, hover effects

### 13.3 Component Structure
```
components/
  ├── Layout.jsx          # Main layout with sidebar
  ├── Modal.jsx           # Reusable modal component
  ├── Button.jsx          # Button variants (if needed)
  ├── Card.jsx            # Card component (if needed)
  └── Table.jsx           # Table component (if needed)
```

### 13.4 State Management
- **Auth**: React Context API
- **Data Fetching**: Axios with interceptors
- **Form State**: React useState hooks
- **Global State**: React Context (if needed)

---

## 14. Design Checklist

### Visual Design
- [ ] All colors meet contrast requirements
- [ ] Typography hierarchy is clear
- [ ] Spacing is consistent throughout
- [ ] Icons are appropriately sized and colored
- [ ] Images/assets are optimized

### Interaction Design
- [ ] All interactive elements have hover states
- [ ] Loading states are clear and informative
- [ ] Error states provide actionable feedback
- [ ] Success states confirm user actions
- [ ] Animations enhance rather than distract

### Responsive Design
- [ ] Mobile layout is tested (< 640px)
- [ ] Tablet layout is tested (640px - 1024px)
- [ ] Desktop layout is tested (> 1024px)
- [ ] Touch targets are at least 44px × 44px
- [ ] Text is readable at all sizes

### Accessibility
- [ ] Color contrast meets WCAG AA standards
- [ ] Keyboard navigation works throughout
- [ ] Screen reader labels are present
- [ ] Focus indicators are visible
- [ ] Motion respects user preferences

---

## 15. Figma Design Structure

### Recommended Figma Organization

```
📁 LogiTrackPro Design
  📁 01_Design System
    📄 Colors
    📄 Typography
    📄 Icons
    📄 Components
  📁 02_Pages
    📄 Login
    📄 Dashboard
    📄 Warehouses
    📄 Customers
    📄 Vehicles
    📄 Plans
    📄 Plan Detail
  📁 03_Components
    📄 Buttons
    📄 Forms
    📄 Cards
    📄 Tables
    📄 Modals
    📄 Navigation
  📁 04_States
    📄 Loading States
    📄 Error States
    📄 Empty States
    📄 Success States
```

### Figma Frame Specifications
- **Desktop**: 1440px × 1024px
- **Tablet**: 768px × 1024px
- **Mobile**: 375px × 812px (iPhone standard)

---

## 16. Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2024 | Initial design specification |

---

## Appendix A: Component Examples

### Button Example
```jsx
<button className="btn btn-primary">
  <Plus className="w-4 h-4" />
  Add New
</button>
```

### Card Example
```jsx
<div className="card card-hover">
  <h3 className="text-lg font-semibold mb-2">Card Title</h3>
  <p className="text-dark-300">Card content goes here</p>
</div>
```

### Badge Example
```jsx
<span className="badge badge-success">Optimized</span>
```

---

**End of Design Specification**

