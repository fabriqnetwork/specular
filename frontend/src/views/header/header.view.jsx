import useHeaderStyles from './header.styles'
import Logo  from '../../images/logo.svg'
import TestnetLogo  from '../../images/testnet-logo.svg'

function Header () {
  const classes = useHeaderStyles()
  return (
    <div className={classes.header}>
                <div>
                    <div>
                        <div className={classes.title}> 
                        <img src={Logo} alt="xDAI to ETH" className={classes.logo}/> Specular Bridge
                        </div>
                    </div>
                </div>
            </div>
  )
}

export default Header
