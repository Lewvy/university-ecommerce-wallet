import { screens as _screens } from 'tailwindcss/defaultTheme';

export const theme = {
	extend: {
		colors: {
			background: 'hsl(var(--background) / <alpha-value>)',
			foreground: 'hsl(var(--foreground) / <alpha-value>)',
			primary: 'hsl(var(--primary) / <alpha-value>)',
			secondary: 'hsl(var(--secondary) / <alpha-value>)',
			border: 'hsl(var(--border) / <alpha-value>)',
			input: 'hsl(var(--input) / <alpha-value>)',
			card: 'hsl(var(--card) / <alpha-value>)',
		},
		borderRadius: {
			lg: `var(--radius)`,
			md: `calc(var(--radius) - 2px)`,
			sm: `calc(var(--radius) - 4px)`,
		},
	},
	screens: {
		..._screens,
	},
};
