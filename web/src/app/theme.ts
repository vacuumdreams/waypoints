import { mergeDeepRight } from 'ramda'

const base = {
  mode: 'light',
  radius: 0,
  transition: 0.2,
  fonts:  {
    size: 14,
    primary: {
      family: "'Nanum Gothic', sans-serif",
      weight: {
        regular: 400,
        bold: 700,
        extrabold: 800,
      },
    },
    script: {
      family: "'Dancing Script', cursive",
      weight: {
        regular: 400,
      },
    },
    title: {
      family: "'Limelight', cursive",
      weight: {
        regular: 400,
      },
    },
  },
  colors: {
    transparent: 'rgba(255,255,255,0)',
    text: {
      main: '#252525',
      weak: '#686868',
    },
    background: {
      main: '#FFFFFF',
      dim: '#F8F8F8',
      contrast: '#121212',
    },
    neutral: {
      main: '#686868',
      weak: '#e2e2e2',
      strong: '#424242',
    },
    primary: {
      main: '#027F99',
      weak: '#029EBD',
      strong: '#236370',
    },
    danger: {
      main: '#FC5735',
    },
    warning:  {
      main: '#FCB203',
    }
  },
}

export const theme = {
  light: base,
  dark: mergeDeepRight(base, {
    mode: 'dark',
    colors: {
      transparent: 'rgba(0,0,0,0)',
      text: {
        main: '#F8F8F8',
      },
      background: {
        main: '#121212',
        dim: '#1E1E1E',
        contrast: '#FFFFFF',
      },
      neutral: {
        weak: '#424242',
        strong: '#e2e2e2',
      },
    },
  }) as typeof base,
}

type ThemCollection = typeof theme;
export type Theme = typeof base;
export type ThemeMode = keyof ThemCollection;
