import { createTheme } from '@mui/material/styles';

export const appTheme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: '#134e4a',
    },
    secondary: {
      main: '#d97706',
    },
    background: {
      default: '#f6f7fb',
      paper: '#ffffff',
    },
  },
  shape: {
    borderRadius: 16,
  },
  typography: {
    fontFamily: ['Inter', 'system-ui', 'sans-serif'].join(','),
  },
});