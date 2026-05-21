---
name: Modern Professional
colors:
  surface: '#f7f9fb'
  surface-dim: '#d8dadc'
  surface-bright: '#f7f9fb'
  surface-container-lowest: '#ffffff'
  surface-container-low: '#f2f4f6'
  surface-container: '#eceef0'
  surface-container-high: '#e6e8ea'
  surface-container-highest: '#e0e3e5'
  on-surface: '#191c1e'
  on-surface-variant: '#45464d'
  inverse-surface: '#2d3133'
  inverse-on-surface: '#eff1f3'
  outline: '#76777d'
  outline-variant: '#c6c6cd'
  surface-tint: '#565e74'
  primary: '#000000'
  on-primary: '#ffffff'
  primary-container: '#131b2e'
  on-primary-container: '#7c839b'
  inverse-primary: '#bec6e0'
  secondary: '#0058be'
  on-secondary: '#ffffff'
  secondary-container: '#2170e4'
  on-secondary-container: '#fefcff'
  tertiary: '#000000'
  on-tertiary: '#ffffff'
  tertiary-container: '#07006c'
  on-tertiary-container: '#7073ff'
  error: '#ba1a1a'
  on-error: '#ffffff'
  error-container: '#ffdad6'
  on-error-container: '#93000a'
  primary-fixed: '#dae2fd'
  primary-fixed-dim: '#bec6e0'
  on-primary-fixed: '#131b2e'
  on-primary-fixed-variant: '#3f465c'
  secondary-fixed: '#d8e2ff'
  secondary-fixed-dim: '#adc6ff'
  on-secondary-fixed: '#001a42'
  on-secondary-fixed-variant: '#004395'
  tertiary-fixed: '#e1e0ff'
  tertiary-fixed-dim: '#c0c1ff'
  on-tertiary-fixed: '#07006c'
  on-tertiary-fixed-variant: '#2f2ebe'
  background: '#f7f9fb'
  on-background: '#191c1e'
  surface-variant: '#e0e3e5'
typography:
  headline-lg:
    fontFamily: Hanken Grotesk
    fontSize: 32px
    fontWeight: '700'
    lineHeight: 40px
    letterSpacing: -0.02em
  headline-lg-mobile:
    fontFamily: Hanken Grotesk
    fontSize: 28px
    fontWeight: '700'
    lineHeight: 36px
    letterSpacing: -0.01em
  headline-md:
    fontFamily: Hanken Grotesk
    fontSize: 24px
    fontWeight: '600'
    lineHeight: 32px
    letterSpacing: -0.01em
  headline-sm:
    fontFamily: Hanken Grotesk
    fontSize: 20px
    fontWeight: '600'
    lineHeight: 28px
  body-lg:
    fontFamily: Hanken Grotesk
    fontSize: 18px
    fontWeight: '400'
    lineHeight: 28px
  body-md:
    fontFamily: Hanken Grotesk
    fontSize: 16px
    fontWeight: '400'
    lineHeight: 24px
  body-sm:
    fontFamily: Hanken Grotesk
    fontSize: 14px
    fontWeight: '400'
    lineHeight: 20px
  label-md:
    fontFamily: Hanken Grotesk
    fontSize: 12px
    fontWeight: '600'
    lineHeight: 16px
    letterSpacing: 0.05em
  label-sm:
    fontFamily: Hanken Grotesk
    fontSize: 11px
    fontWeight: '500'
    lineHeight: 14px
rounded:
  sm: 0.25rem
  DEFAULT: 0.5rem
  md: 0.75rem
  lg: 1rem
  xl: 1.5rem
  full: 9999px
spacing:
  base: 8px
  xs: 4px
  sm: 12px
  md: 24px
  lg: 40px
  xl: 64px
  gutter: 24px
  margin-mobile: 16px
  margin-desktop: 48px
---

## Brand & Style
The design system is anchored in a **Corporate / Modern** aesthetic, prioritizing clarity, efficiency, and a high-degree of perceived reliability. It is engineered for professional environments where information density must coexist with visual breathing room. 

The visual language balances a neutral, systematic foundation with subtle premium touches. By utilizing a minimalist approach to ornamentation and focusing on precise alignments and purposeful whitespace, the UI evokes a sense of calm authority. The goal is to create a tool-like experience that feels sophisticated yet stays out of the user's way, allowing data and content to remain the focal point.

## Colors
The palette is built on a "Deep Slate" primary foundation to provide a sense of stability and professional gravity. The secondary "Electric Blue" is reserved for primary actions and interactive states, ensuring high discoverability without overwhelming the interface.

*   **Primary:** Used for text, iconography, and high-emphasis structural elements.
*   **Secondary:** The interactive accent for buttons, links, and active states.
*   **Tertiary:** Used sparingly for data visualization or secondary highlights to provide depth.
*   **Neutral:** A range of cool grays used for backgrounds, borders, and subtle containment to maintain a clean, organized hierarchy.

## Typography
This design system utilizes **Hanken Grotesk** across all levels to achieve a contemporary, precise, and highly legible typographic character. The typeface’s sharp geometry and modern proportions provide the "Matter-like" aesthetic—clean, professional, and sophisticated.

Headlines use tighter letter-spacing and heavier weights to create a strong visual anchor. Body copy prioritizes readability with generous line heights. Labels are set with increased tracking and semi-bold weights to ensure clarity at small scales, particularly in data-heavy views.

## Layout & Spacing
The layout follows a **Fluid Grid** philosophy based on an 8px base unit. This ensures mathematical harmony across all components and screen sizes.

*   **Desktop:** A 12-column grid with 24px gutters and 48px side margins. Content typically spans 6, 8, or 12 columns.
*   **Tablet:** An 8-column grid with 16px gutters and 24px side margins.
*   **Mobile:** A 4-column grid with 16px gutters and 16px side margins.

Spacing is applied through a strict scale to maintain vertical rhythm. Use `md` (24px) for most component grouping and `lg` (40px) for section headers.

## Elevation & Depth
Depth is communicated through **Tonal Layers** and **Ambient Shadows**. The system avoids heavy gradients, opting instead for subtle surface-container tiers to define hierarchy.

*   **Level 0 (Base):** The main background using the Neutral color.
*   **Level 1 (Surface):** White or slightly off-white surfaces for cards and main containers, featuring a very soft, diffused shadow (0px 4px 20px rgba(0,0,0,0.04)).
*   **Level 2 (Overlay):** Used for modals and dropdowns, utilizing a more pronounced ambient shadow (0px 12px 32px rgba(0,0,0,0.08)) to clear the base content.
*   **Interactive State:** Subtle lift on hover (0px 8px 24px rgba(0,0,0,0.06)) to provide tactile feedback.

## Shapes
The shape language is **Rounded**, utilizing a 0.5rem (8px) corner radius as the standard for buttons, inputs, and cards. This radius balances the precision of the typography with a touch of approachability. Large containers and sections may scale up to `rounded-xl` (1.5rem) to soften the overall layout, while small elements like tags or badges use a fully pill-shaped radius for distinct visual separation.

## Components
Consistent component styling is vital for the integrity of the design system:

*   **Buttons:** Primary buttons feature solid Primary fills with Hanken Grotesk SemiBold text. Secondary buttons use a subtle gray stroke with no fill.
*   **Input Fields:** Border-based with a 1px stroke in a light neutral. On focus, the border transitions to the Secondary color with a soft 2px outer glow.
*   **Cards:** White backgrounds, Level 1 elevation shadows, and 8px (roundedness-2) corners. Padding should be consistent at `md` (24px).
*   **Chips/Tags:** Small, low-contrast background fills with `label-sm` typography and pill-shaped corners.
*   **Lists:** High-density with 1px horizontal dividers. Interactive list items should feature a subtle background color shift on hover.
*   **Checkboxes/Radios:** Use the Secondary color for active states, maintaining sharp geometric clarity to match the typeface.