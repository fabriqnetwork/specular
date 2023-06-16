import React from 'react'
import ReactDOM from 'react-dom'
import { ThemeProvider } from 'react-jss'

import reportWebVitals from './report-web-vitals'
import theme from './styles/theme'
import App from './views/app.view'

ReactDOM.render(
  <React.StrictMode>
    <ThemeProvider theme={theme}>
      <App />
    </ThemeProvider>
  </React.StrictMode>,
  document.getElementById('root')
)

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals()
