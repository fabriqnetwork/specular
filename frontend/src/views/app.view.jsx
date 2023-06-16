import useAppStyles from './app.styles'
import Stepper from './stepper/stepper.view'
import Footer from './footer/footer.view'
import Header from './header/header.view'
import FAQ from './faq/faq.view'
import { useState } from 'react'



function App () {
  const classes=useAppStyles();
  const [setOpenGetMoreFaq] = useState(() => () => null);

  return <div className={classes.container}> 
    <Header/>
    <Stepper/>
    <FAQ setOpenGetMoreFaq={setOpenGetMoreFaq} />
    <Footer/>
  </div>

}

export default App
