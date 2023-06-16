import useFooterStyles from './footer.styles'
import {Container } from "@mui/system"
import twitter from '../../images/twitter.svg';
import github from '../../images/github.svg';
import medium from '../../images/medium.svg';
import { Grid, Typography } from "@mui/material";
import Logo  from '../../images/footer-logo.svg'
import TestnetLogo  from '../../images/footer-logo-testnet.svg'
import useWindowSize from '../../hooks/use-window-size'

const footerLinks = [
  {
      title: "Docs",
      links: [
          {
              title: "Getting Started",
              url: "https://specular.network/docs/getting-started",
          },
          {
              title: "Community",
              url: "https://specular.network/docs/community/",
          }
      ],
  },
];
function Footer () {
  const classes = useFooterStyles()
  const scrollToTop = () => {
    window.scrollTo({ top: 0, left: 0, behavior: "smooth" });
};

let newDate = new Date()
let date = newDate.getDate();
let month = newDate.getMonth() + 1;
let year = newDate.getFullYear();
var pjson = require('../../../package.json');
let version = pjson.version;
let useProd; 
if (process.env.REACT_APP_PROD === 'true') {
    useProd = <></>;
  } else {
    useProd = <div className={classes.bottom}>
                <Typography
                fontSize="10px"
                color="grey"
                >
                Version {version} and Date {date}/{month}/{year}
                </Typography>
            </div>;
  }
const size = useWindowSize(); 

let footerLogo;

if (size.width > 675){
    footerLogo = <Grid item onClick={scrollToTop}>
                    <img src={Logo} alt="xDAI to ETH" className={classes.logo}/>
                </Grid>  
} else {
    footerLogo = <></>
}
  return (
    <footer className={classes.footer}>
            <Container>
                    <div>
                        <Grid
                            container
                            spacing={10}
                            flex="1"
                            justifyContent="center"
                        >
                            {footerLogo}
                            {footerLinks.map((footer) => (
                                <Grid item key={footer.title}>
                                    <Typography
                                        variant="h6"
                                        fontWeight="bold"
                                        fontSize="14px"
                                        color="white"
                                        align="center"
                                        sx={{
                                            textTransform: "uppercase",
                                            marginBottom: "10px",
                                        }}
                                    >
                                        {footer.title}
                                    </Typography>
                                    {footer.links.map((link) => (
                                        <Typography
                                            key={link.title}
                                            variant="subtitle1"
                                            fontSize="14px"
                                            color="grey"
                                        >
                                            <a
                                                target="_blank"
                                                rel="noreferrer noopener"
                                                href={link.url}
                                            >
                                                {link.title}
                                            </a>
                                        </Typography>
                                    ))}
                                </Grid>
                            ))}
                            <Grid item>
                                <Typography
                                    variant="h6"
                                    fontWeight="bold"
                                    fontSize="14px"
                                    color="white"
                                    align="center"
                                    sx={{
                                        textTransform: "uppercase",
                                        marginBottom: "15px",
                                    }}
                                >
                                    Socials
                                </Typography>
                                <Grid container spacing={2}>
                                    <Grid item>
                                        <div>
                                            <a
                                                href="https://twitter.com/specularl2"
                                                target="_blank"
                                                rel="noreferrer noopener"
                                            >
                                                <img
                                                    src={twitter}
                                                    alt="Specular Twitter"
                                                />
                                            </a>
                                        </div>
                                    </Grid>
                                    <Grid item>
                                        <div>
                                            <a
                                                href="https://github.com/specularl2"
                                                target="_blank"
                                                rel="noreferrer noopener"
                                            >
                                                <img
                                                    src={github}
                                                    alt="Specular Github"
                                                />
                                            </a>
                                        </div>
                                    </Grid>
                                    <Grid item>
                                        <div>
                                            <a
                                                href="https://medium.com/@SpecularL2"
                                                target="_blank"
                                                rel="noreferrer noopener"
                                            >
                                                <img
                                                    src={medium}
                                                    alt="Specular Medium"
                                                />
                                            </a>
                                        </div>
                                    </Grid>
                                </Grid>
                            </Grid>
                        </Grid>
                    </div>
            </Container>
            {useProd}
        </footer>
  )
}

export default Footer
