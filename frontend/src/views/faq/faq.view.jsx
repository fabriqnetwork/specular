import {
    AccordionDetails,
    AccordionSummary,
    Typography,
} from "@mui/material";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import { Container } from "@mui/system";
import { useState, useEffect, useRef } from "react";
import { styled } from '@mui/material/styles';
import MuiAccordion from '@mui/material/Accordion';


const FAQ = ({ setOpenGetMoreFaq }) => {
    const [expanded, setExpanded] = useState("panel1");
    const handleChange = (panel) => (event, newExpanded) => {
        setExpanded(newExpanded ? panel : false);
    };

    const getMoreFaqRef = useRef(null);

    useEffect(() => {
        setOpenGetMoreFaq(() => () => {
            setExpanded("hcigmx");
            if (getMoreFaqRef.current !== null) {
                getMoreFaqRef.current.scrollIntoView({
                    behavior: "smooth",
                    block: "center",
                    inline: "center",
                });
            }
        });
    }, []);

    const Accordion = styled((props) => (
        <MuiAccordion disableGutters elevation={0} square {...props} />
      ))(({ theme }) => ({
        border: `1px solid ${theme.palette.divider}`,
        borderRadius:'16px',
        marginTop: "0.32em",
    
        '&:before': {
          display: 'none',
        },
      }));

    return (
        <Container maxWidth="sm">
            <Typography
                align="center"
                marginTop="1em"
                fontSize="32px"
            >
                FAQ
            </Typography>


            <Accordion
                expanded={expanded === "one"}
                onChange={handleChange("one")}
            >
                <AccordionSummary
                    expandIcon={<ExpandMoreIcon />}
                    aria-controls="panel2a-content"
                    id="panel2a-header"
                >
                    <Typography
                        color="Grey"
                        fontSize="18px"
                    >
                        What is Specular?
                    </Typography>
                </AccordionSummary>
                <AccordionDetails>
                    <Typography
                        align="justify"
                        fontSize="16px"
                    >
                        <span>
                        Specular Network is a Layer 2 scaling solution for Ethereum that uses optimistic rollup technology to scale the Ethereum network. It is still under development, but it has the potential to provide a number of benefits for Ethereum users, including scalability, security, and decentralization.
                        </span>
                    </Typography>
                </AccordionDetails>
            </Accordion>

            <Accordion
                expanded={expanded === "two"}
                onChange={handleChange("two")}
            >
                <AccordionSummary
                    expandIcon={<ExpandMoreIcon />}
                    aria-controls="panel2a-content"
                    id="panel2a-header"
                >
                    <Typography
                        color="Grey"
                        fontSize="18px"
                    >
                        What is Specular Bridge?
                    </Typography>
                </AccordionSummary>
                <AccordionDetails>
                    <Typography
                        align="justify"
                        fontSize="16px"
                    >
                        <span>
                        The Specular Network Bridge is a two-way bridge that allows users to transfer assets between Ethereum and Specular Network. It is still in beta, but it is expected to be released to the public in the near future. The bridge offers low fees, high throughput, and security.
                        </span>
                    </Typography>
                </AccordionDetails>
            </Accordion>

            <Accordion
                expanded={expanded === "three"}
                onChange={handleChange("three")}
            >
                <AccordionSummary
                    expandIcon={<ExpandMoreIcon />}
                    aria-controls="panel2a-content"
                    id="panel2a-header"
                >
                    <Typography
                        color="Grey"
                        fontSize="18px"
                    >
                        What chains are supported By Specular Bridge?
                    </Typography>
                </AccordionSummary>
                <AccordionDetails>
                    <Typography
                        align="justify"
                        fontSize="16px"
                    >
                        <span>
                        As of today, Specular Network only supports Chiado. However, the team has plans to support other chains in the future.
                        </span>
                    </Typography>
                </AccordionDetails>
            </Accordion>

            <div ref={getMoreFaqRef} />
        </Container>
    );
};

export default FAQ;

