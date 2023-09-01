import useHeaderStyles from './header.styles'
import Logo  from '../../images/logo.svg'

function Header () {
  const classes = useHeaderStyles()
  return (
    <div className={classes.header}>
                <div>
                    <div>
                        <div className={classes.title}>
                        <img src={Logo} alt="Specular Bridge" className={classes.logo}/>
                        </div>
                    </div>
                </div>
            </div>
  )
}

export default Header
